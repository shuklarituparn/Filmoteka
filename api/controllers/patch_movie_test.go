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

func TestPatchMovie(t *testing.T) {
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
		return
	}
	patchData := models.Movie{
		Title:       "Title is patched",
		Description: "Desription is patched",
	}
	patchMovieJSON, err := json.Marshal(patchData)
	if err != nil {
		t.Error(err)
	}
	patchMovieReq, err := http.NewRequest("PATCH", "/api/v1/movies/patch?id="+strconv.Itoa(int(movieResponse.Data.ID)), bytes.NewBuffer(patchMovieJSON))
	if err != nil {
		t.Error(err)
	}
	patchMovieReq.Header.Set("Authorization", "Bearer "+accessToken)
	patchActorRecorder := httptest.NewRecorder()
	PatchMovie(db)(patchActorRecorder, patchMovieReq)
	if patchActorRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, patchActorRecorder.Code)
	}
	defer CloseAndDelete(db)
}
