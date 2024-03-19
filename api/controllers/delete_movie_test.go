package controllers

import (
	"bytes"
	"encoding/json"

	"github.com/shuklarituparn/Filmoteka/api/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestDeleteMovie(t *testing.T) {
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
	CreateActor(db)(createRecorder, createMovie)
	var movieResponse CreateMovieResponse
	err = json.NewDecoder(createRecorder.Body).Decode(&movieResponse)
	if err != nil {
		return
	}
	deleteReq, err := http.NewRequest("DELETE", "/api/v1/movies/delete?id="+strconv.Itoa(int(movieResponse.Data.ID)), bytes.NewReader([]byte{}))
	if err != nil {
		t.Error(err)
	}
	deleteReq.Header.Set("Authorization", "Bearer "+accessToken)
	deleteRecorder := httptest.NewRecorder()
	DeleteMovie(db)(deleteRecorder, deleteReq)
	if deleteRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, deleteRecorder.Code)
	}
	defer CloseAndDelete(db)
}
