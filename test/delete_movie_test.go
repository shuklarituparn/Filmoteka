package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestDeleteMovie(t *testing.T) {

	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking movie deletion", func(P provider.T) {
		P.Feature("Movie Deletion")
		P.Title("Testing movie deletion")
		P.Description("This test checks if the movie can be deleted or not")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		movie := models.CreateMovieModel{
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
		controllers.CreateMovie(db)(createRecorder, createMovie)
		var movieResponse controllers.CreateMovieResponse
		err = json.NewDecoder(createRecorder.Body).Decode(&movieResponse)
		if err != nil {
			return
		}
		deleteReq, err := http.NewRequest("DELETE", "/api/v1/movies/delete?id="+strconv.Itoa(int(movieResponse.Data.ID)), nil)
		if err != nil {
			t.Error(err)
		}
		deleteReq.Header.Set("Authorization", "Bearer "+accessToken)
		deleteRecorder := httptest.NewRecorder()
		controllers.DeleteMovie(db)(deleteRecorder, deleteReq)
		if deleteRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, deleteRecorder.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}

func TestDeleteMovieWithData(t *testing.T) {

	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking movie deletion", func(P provider.T) {
		P.Feature("Movie Deletion")
		P.Title("Testing movie deletion without a body")
		P.Description("This test checks if the movie can be deleted without a body or not")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		movie := models.CreateMovieModel{
			Title:  "",
			Actors: nil,
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
		controllers.CreateMovie(db)(createRecorder, createMovie)
		var movieResponse controllers.CreateMovieResponse
		fmt.Println(movieResponse.Data.ID)
		err = json.NewDecoder(createRecorder.Body).Decode(&movieResponse)
		if err != nil {
			return
		}
		deleteReq, err := http.NewRequest("DELETE", "/api/v1/movies/delete?id="+strconv.Itoa(int(movieResponse.Data.ID)), nil)
		if err != nil {
			t.Error(err)
		}
		deleteReq.Header.Set("Authorization", "Bearer "+accessToken)
		deleteRecorder := httptest.NewRecorder()
		controllers.DeleteMovie(db)(deleteRecorder, deleteReq)
		if deleteRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, deleteRecorder.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}
