package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"movie-list/logger"
	"os"

	"github.com/go-sql-driver/mysql"
)

type DBMovie struct {
	Id          int64
	Title       string
	Year        string
	ImdbId      string
	MALId       int
	ShowType    string
	Poster      string
	UserScore   float64
	Plot        string
	Director    string
	Language    string
	Genre       string
	Released    string
	Runtime     string
	Metascore   string
	ImdbRating  string
	Watched     bool
	ListType    string
	DateWatched string
	Notes       string
}

type DBObject struct {
	isJson        bool
	db            *sql.DB
	jsonPath      string
	encrypted     bool
	encryptionKey []byte
}

var dbObject DBObject

func AddMovie(movie DBMovie) (int64, error) {
	if dbObject.isJson {
		movies, err := dbObject.getMoviesJson()
		if err != nil {
			return 0, fmt.Errorf("add_movie: %v", err)
		}
		movies = append(movies, movie)
		filedata, err := json.MarshalIndent(movies, "", "  ")
		if err != nil {
			return 0, fmt.Errorf("add_movie: %v", err)
		}
		if dbObject.encrypted {
			filedata, err = encrypt(dbObject.encryptionKey, filedata)
			if err != nil {
				return 0, fmt.Errorf("add_movie: %v", err)
			}
		}
		// Write the JSON data to the file
		os.WriteFile(dbObject.jsonPath, filedata, 0644)
		return int64(len(movies)), nil
	} else {
		return dbObject.addMovieDb(movie)
	}
}
func GetMovies() ([]DBMovie, error) {
	if dbObject.isJson {
		return dbObject.getMoviesJson()
	} else {
		return dbObject.getMoviesDb()
	}

}
func ConnectDBObject(isJson bool, isEncrypted bool) {
	if isJson {
		dbObject = connectDbJson(isEncrypted)
	} else {
		dbObject = connectDb()
	}
}
func connectDbJson(encrypted bool) DBObject {
	d := DBObject{
		isJson:    true,
		db:        nil,
		jsonPath:  os.Getenv("JSON_PATH"),
		encrypted: encrypted,
	}
	if encrypted {
		d.encryptionKey = deriveKey(os.Getenv("ENCRYPTION_KEY"))
	}
	return d
}

func connectDb() DBObject {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_HOST"),
		DBName:               "go_movies",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Log.Fatalf("sql error: %s", err.Error())
	}

	pingErr := db.Ping()
	if pingErr != nil {
		logger.Log.Fatalf("sql error: %s", err.Error())
	}
	fmt.Println("Connected!")
	return DBObject{
		isJson:   false,
		db:       db,
		jsonPath: "",
	}
}
func (d *DBObject) addMovieDb(movie DBMovie) (int64, error) {
	dateWatched := sql.NullString{}
	if movie.DateWatched != "" {
		dateWatched.String = movie.DateWatched
		dateWatched.Valid = true
	} else {
		dateWatched.Valid = false
	}
	query := "INSERT INTO movies (Title, Year, ImdbId, MALId, ShowType, Poster, UserScore, Plot, Director, Language, Genre, Released, Runtime, Metascore, ImdbRating,Watched, ListType, DateWatched, Notes) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := d.db.Exec(query, movie.Title, movie.Year, movie.ImdbId, movie.MALId, movie.ShowType, movie.Poster, movie.UserScore,
		movie.Plot, movie.Director, movie.Language, movie.Genre, movie.Released, movie.Runtime, movie.Metascore, movie.ImdbRating, movie.Watched, movie.ListType, dateWatched, movie.Notes)
	if err != nil {
		return 0, fmt.Errorf("add_movie: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add_movie: %v", err)
	}
	return id, nil
}
func (d *DBObject) getMoviesDb() ([]DBMovie, error) {
	query := "SELECT * FROM movies"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get_movies: %v", err)
	}
	defer rows.Close()

	movies := []DBMovie{}
	for rows.Next() {
		var movie DBMovie
		var dateWatched sql.NullString
		err := rows.Scan(&movie.Id, &movie.Title, &movie.Year, &movie.ImdbId, &movie.MALId, &movie.ShowType, &movie.Poster,
			&movie.UserScore, &movie.Plot, &movie.Director, &movie.Language, &movie.Genre, &movie.Released,
			&movie.Runtime, &movie.Metascore, &movie.ImdbRating, &movie.Watched, &movie.ListType, &dateWatched, &movie.Notes)
		if err != nil {
			return nil, fmt.Errorf("get_movies: %v", err)
		}
		if dateWatched.Valid {
			movie.DateWatched = dateWatched.String
		} else {
			movie.DateWatched = ""
		}
		movies = append(movies, movie)
	}
	return movies, nil
}
func (d *DBObject) getMoviesJson() ([]DBMovie, error) {
	filename := d.jsonPath
	//check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		file, err := os.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("get_movies: %v", err)
		}
		defer file.Close()
		filedata, err := json.MarshalIndent([]DBMovie{}, "", "  ")
		if d.encrypted {
			filedata, err = encrypt(d.encryptionKey, filedata)
		}
		os.WriteFile(filename, filedata, 0644)
		if err != nil {
			return nil, fmt.Errorf("get_movies: %v", err)
		}
		return []DBMovie{}, nil
	}
	filedata, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("get_movies: %v", err)
	}

	if d.encrypted {
		filedata, err = decrypt(d.encryptionKey, filedata)
		if err != nil {
			return nil, fmt.Errorf("get_movies: %v", err)
		}
	}

	var movies []DBMovie
	err = json.Unmarshal(filedata, &movies)
	if err != nil {
		return nil, fmt.Errorf("get_movies: %v", err)
	}
	return movies, nil
}
func (d *DBObject) Close() error {
	if d.isJson {
		return nil
	} else {
		return d.db.Close()
	}
}
