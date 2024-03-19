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

func TestReadAMovie(t *testing.T) {
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
	createReq, err := http.NewRequest("POST", "/api/v1/movies/create", bytes.NewBuffer(movieJSON))
	if err != nil {
		t.Error(err)
	}
	createReq.Header.Set("Authorization", "Bearer "+accessToken)
	createRecorder := httptest.NewRecorder()
	CreateMovie(db)(createRecorder, createReq)
	if createRecorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusOK, createRecorder.Code)
	}
	var movieResponse CreateMovieResponse
	err = json.NewDecoder(createRecorder.Body).Decode(&movieResponse)
	if err != nil {
		return
	}
	getReq, err := http.NewRequest("GET", "/api/v1/movies/get?id="+strconv.Itoa(int(movieResponse.Data.ID)), bytes.NewReader([]byte{}))
	if err != nil {
		t.Error(err)
	}
	getReq.Header.Set("Authorization", "Bearer "+accessToken)
	getRecorder := httptest.NewRecorder()
	ReadMovie(db)(getRecorder, getReq)
	if getRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, getRecorder.Code)
	}
	defer CloseAndDelete(db)
}
