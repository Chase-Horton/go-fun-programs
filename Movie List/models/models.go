package models

import (
	"encoding/json"
	"fmt"
	"io"
	"movie-list/logger"
	"net/http"
	"os"
)

type MovieResponse struct {
	Movies       []Movie `json:"search"`
	TotalResults string  `json:"totalResults"`
	Response     string  `json:"response"`
}
type Movie struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	ImdbId   string `json:"imdbId"`
	ShowType string `json:"Type"`
	Poster   string `json:"Poster"`
}
type MovieDetailed struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	Released   string `json:"Released"`
	Rated      string `json:"Rated"`
	Runtime    string `json:"Runtime"`
	Genre      string `json:"Genre"`
	Director   string `json:"Director"`
	Plot       string `json:"Plot"`
	Language   string `json:"Language"`
	Poster     string `json:"Poster"`
	ImdbRating string `json:"imdbRating"`
	Metascore  string `json:"Metascore"`
	ImdbId     string `json:"imdbId"`
	ShowType   string `json:"Type"`
}

func (m MovieDetailed) ToDBMovie(score float64) DBMovie {
	return DBMovie{
		Title:      m.Title,
		Year:       m.Year,
		ImdbId:     m.ImdbId,
		MALId:      -1,
		ShowType:   m.ShowType,
		Poster:     m.Poster,
		UserScore:  score,
		Plot:       m.Plot,
		Director:   m.Director,
		Language:   m.Language,
		Genre:      m.Genre,
		Released:   m.Released,
		Runtime:    m.Runtime,
		Metascore:  m.Metascore,
		ImdbRating: m.ImdbRating,
		Watched:    true,
		ListType:   "watched"}
}

func (m Movie) String() string {
	return fmt.Sprintf("%s (%s)", m.Title, m.Year)
}

var apiKey string

func GetApiKey() string {
	if apiKey != "" {
		return apiKey
	}
	apiKey = os.Getenv("OMDB_API_KEY")
	return apiKey
}
func SearchMovie(title string) MovieResponse {
	response, err := http.Get("http://www.omdbapi.com/?apikey=" + GetApiKey() + "&s=" + title)
	if err != nil {
		logger.Log.Fatalf("Error Searching for Movie: %s", err.Error())
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Log.Fatalf("Error Parsing Search Movie Response: %s", err.Error())
	}

	var responseObject MovieResponse
	json.Unmarshal(responseData, &responseObject)
	return responseObject
}
func GetMovieDetails(title string) MovieDetailed {
	response, err := http.Get("http://www.omdbapi.com/?apikey=" + GetApiKey() + "&i=" + title + "&plot=full")
	if err != nil {
		logger.Log.Fatalf("Error Getting Movie Details: %s", err.Error())
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Log.Fatalf("Error Parsing Movie Details: %s", err.Error())
	}

	var responseObject MovieDetailed
	json.Unmarshal(responseData, &responseObject)
	return responseObject
}
