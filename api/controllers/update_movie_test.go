package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMovie(t *testing.T) {
	ConnectTestDB()
	var db = GetDBInstance()
	accessToken := LoginAndGetAccessToken(t, db)
	movie := models.Movie{
		Description: "This is my movie Descripton",
		Genre:       "Drama",
		Rating:      6.7,
		ReleaseYear: 2000,
		Title:       "Movie Title",
	}
	movieJSON, err := json.Marshal(movie)
	if err != nil {
		t.Error(err)
	}
	createMovie, err := http.NewRequest("POST", "/api/v1/movies/create", bytes.NewBuffer(movieJSON))
	if err != nil {
		t.Error(err)
	}
	createMovie.Header.Set("Authorization", "Bearer "+accessToken)
	createRecorder := httptest.NewRecorder()
	CreateMovie(db)(createRecorder, createMovie)
	var movieResponse CreateMovieResponse
	err = json.NewDecoder(createRecorder.Body).Decode(&movieResponse)
	if err != nil {
		t.Error(err)
	}
	updateData := models.Movie{
		ID:          movieResponse.Data.ID,
		Description: "This is my movie Updated",
		Genre:       "Drama Updated",
		Rating:      3.7,
		ReleaseYear: 1000,
		Title:       "Movie Title Updated",
	}
	updateMovieJSON, err := json.Marshal(updateData)
	if err != nil {
		t.Error(err)
	}
	updateMovieReq, err := http.NewRequest("PUT", "/api/v1/movies/update", bytes.NewBuffer(updateMovieJSON))
	if err != nil {
		t.Error(err)
	}
	updateMovieReq.Header.Set("Authorization", "Bearer "+accessToken)
	updateMovieRecorder := httptest.NewRecorder()
	UpdateMovie(db)(updateMovieRecorder, updateMovieReq)
	if updateMovieRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, updateMovieRecorder.Code)
	}
	defer CloseAndDelete(db)
}
