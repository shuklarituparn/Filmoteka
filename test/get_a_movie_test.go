package test

import (
	"bytes"
	"encoding/json"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestReadAMovie(t *testing.T) {
	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking movie read", func(P provider.T) {
		P.Feature("Movie getting")
		P.Title("Testing movie read")
		P.Description("This test checks if the movie can be read or not")
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
		controllers.CreateMovie(db)(createRecorder, createReq)
		if createRecorder.Code != http.StatusCreated {
			t.Errorf("Expected status %d; got %d", http.StatusOK, createRecorder.Code)
		}
		var movieResponse controllers.CreateMovieResponse
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
		controllers.ReadMovie(db)(getRecorder, getReq)
		if getRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, getRecorder.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}

func TestReadAMovieWithoutBody(t *testing.T) {
	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking movie read", func(P provider.T) {
		P.Feature("Movie getting")
		P.Title("Testing movie read")
		P.Description("This test checks if the movie can be read or not with a broken link")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		movie := models.CreateMovieModel{}
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
		controllers.CreateMovie(db)(createRecorder, createReq)
		if createRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, createRecorder.Code)
		}
		var movieResponse controllers.CreateMovieResponse
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
		controllers.ReadMovie(db)(getRecorder, getReq)
		if getRecorder.Code != http.StatusNotFound {
			t.Errorf("Expected status %d; got %d", http.StatusNotFound, getRecorder.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}
