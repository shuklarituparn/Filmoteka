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
	"testing"
)

func TestCreateMovie(t *testing.T) {

	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking movie creation", func(P provider.T) {
		P.Feature("Movie Creation")
		P.Title("Testing movie creation")
		P.Description("This test checks if the movie can be created or not")
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
		defer CloseAndDelete(db)

	})
	r.RunTests()

}

func TestCreateMovieWithoutBody(t *testing.T) {

	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking false movie creation", func(P provider.T) {
		P.Feature("Movie Creation")
		P.Title("Testing movie creation")
		P.Description("This test checks if the movie can be created or not without a body")
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
		defer CloseAndDelete(db)

	})
	r.RunTests()

}
