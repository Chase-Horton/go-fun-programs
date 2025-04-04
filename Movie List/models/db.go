package models

import (
	"database/sql"
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

var db *sql.DB

func ConnectDb() {
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
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Log.Fatalf("sql error: %s", err.Error())
	}

	pingErr := db.Ping()
	if pingErr != nil {
		logger.Log.Fatalf("sql error: %s", err.Error())
	}
	fmt.Println("Connected!")
}
func AddMovie(movie DBMovie) (int64, error) {
	query := "INSERT INTO movies (Title, Year, ImdbId, MALId, ShowType, Poster, UserScore, Plot, Director, Language, Genre, Released, Runtime, Metascore, ImdbRating,Watched, ListType, DateWatched, Notes) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, movie.Title, movie.Year, movie.ImdbId, movie.MALId, movie.ShowType, movie.Poster, movie.UserScore,
		movie.Plot, movie.Director, movie.Language, movie.Genre, movie.Released, movie.Runtime, movie.Metascore, movie.ImdbRating, movie.Watched, movie.ListType, movie.DateWatched, movie.Notes)
	if err != nil {
		return 0, fmt.Errorf("add_movie: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add_movie: %v", err)
	}
	return id, nil
}
func GetMovies() ([]DBMovie, error) {
	query := "SELECT * FROM movies"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get_movies: %v", err)
	}
	defer rows.Close()

	movies := []DBMovie{}
	for rows.Next() {
		var movie DBMovie
		err := rows.Scan(&movie.Id, &movie.Title, &movie.Year, &movie.ImdbId, &movie.MALId, &movie.ShowType, &movie.Poster,
			&movie.UserScore, &movie.Plot, &movie.Director, &movie.Language, &movie.Genre, &movie.Released,
			&movie.Runtime, &movie.Metascore, &movie.ImdbRating, &movie.Watched, &movie.ListType, &movie.DateWatched, &movie.Notes)
		if err != nil {
			return nil, fmt.Errorf("get_movies: %v", err)
		}
		movies = append(movies, movie)
	}
	return movies, nil
}
