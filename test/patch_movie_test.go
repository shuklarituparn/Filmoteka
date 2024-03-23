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

func TestPatchMovie(t *testing.T) {
	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking Movie Patch", func(P provider.T) {
		P.Feature("Movie Patch")
		P.Title("Testing movie patch")
		P.Description("This test checks if the movie can be patched or not")
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
		controllers.CreateMovie(db)(createRecorder, createMovie)
		var movieResponse controllers.CreateMovieResponse
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
		controllers.PatchMovie(db)(patchActorRecorder, patchMovieReq)
		if patchActorRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, patchActorRecorder.Code)
		}
		defer CloseAndDelete(db)

	})
	r.RunTests()

}

func TestPatchMovieWithoutBody(t *testing.T) {
	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking Movie Patch", func(P provider.T) {
		P.Feature("Movie Patch")
		P.Title("Testing movie patch without body")
		P.Description("This test checks if the movie can be patched or not without body")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		movie := models.Movie{}
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
		patchMovieJSON := []byte{}
		patchMovieReq, err := http.NewRequest("PATCH", "/api/v1/movies/patch?id="+strconv.Itoa(int(movieResponse.Data.ID)), bytes.NewBuffer(patchMovieJSON))
		if err != nil {
			t.Error(err)
		}
		patchMovieReq.Header.Set("Authorization", "Bearer "+accessToken)
		patchActorRecorder := httptest.NewRecorder()
		controllers.PatchMovie(db)(patchActorRecorder, patchMovieReq)
		if patchActorRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, patchActorRecorder.Code)
		}
		defer CloseAndDelete(db)

	})
	r.RunTests()

}

func TestPatchMovieWithEmptyBody(t *testing.T) {
	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking movie Patch", func(P provider.T) {
		P.Feature("Movie Patch")
		P.Title("Testing movie patch without body")
		P.Description("This test checks if the movie can be patched or not without body")
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
		controllers.CreateMovie(db)(createRecorder, createMovie)
		var movieResponse controllers.CreateMovieResponse
		err = json.NewDecoder(createRecorder.Body).Decode(&movieResponse)
		if err != nil {
			return
		}
		patchMovieJSON := []byte{}
		patchMovieReq, err := http.NewRequest("PATCH", "/api/v1/movies/patch?id="+strconv.Itoa(int(movieResponse.Data.ID)), bytes.NewBuffer(patchMovieJSON))
		if err != nil {
			t.Error(err)
		}
		patchMovieReq.Header.Set("Authorization", "Bearer "+accessToken)
		patchActorRecorder := httptest.NewRecorder()
		controllers.PatchMovie(db)(patchActorRecorder, patchMovieReq)
		if patchActorRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, patchActorRecorder.Code)
		}
		defer CloseAndDelete(db)

	})
	r.RunTests()

}

func TestPatchMovieWithInvalidID(t *testing.T) {
	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking Movie Patch", func(P provider.T) {
		P.Feature("Movie Patch")
		P.Title("Testing movie patch without body")
		P.Description("This test checks if the movie can be patched or not with an invalid id")
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
		controllers.CreateMovie(db)(createRecorder, createMovie)
		var movieResponse controllers.CreateMovieResponse
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
		patchMovieReq, err := http.NewRequest("PATCH", "/api/v1/movies/patch", bytes.NewBuffer(patchMovieJSON))
		if err != nil {
			t.Error(err)
		}
		patchMovieReq.Header.Set("Authorization", "Bearer "+accessToken)
		patchActorRecorder := httptest.NewRecorder()
		controllers.PatchMovie(db)(patchActorRecorder, patchMovieReq)
		if patchActorRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, patchActorRecorder.Code)
		}
		defer CloseAndDelete(db)

	})
	r.RunTests()

}
