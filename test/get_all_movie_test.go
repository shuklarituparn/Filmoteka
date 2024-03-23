package test

import (
	"bytes"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadAllMovies(t *testing.T) {
	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking movie read ", func(P provider.T) {
		P.Feature("Movie getting")
		P.Title("Testing movie read")
		P.Description("This test checks if all the movie can be read")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		getReq, err := http.NewRequest("GET", "/api/v1/movies/all?page_size=3&page=1", bytes.NewReader([]byte{}))
		if err != nil {
			t.Error(err)
		}
		getReq.Header.Set("Authorization", "Bearer "+accessToken)
		getRecorder := httptest.NewRecorder()
		controllers.ReadAllMovies(db)(getRecorder, getReq)
		if getRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, getRecorder.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}

func TestReadAllMoviesWithoutPagination(t *testing.T) {
	r := runner.NewRunner(t, "Movie")
	r.NewTest("Checking movie read without pagination", func(P provider.T) {
		P.Feature("Movie getting")
		P.Title("Testing movie read")
		P.Description("This test checks if the movie can be read or not without pagination")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		getReq, err := http.NewRequest("GET", "/api/v1/movies/all", bytes.NewReader([]byte{}))
		if err != nil {
			t.Error(err)
		}
		getReq.Header.Set("Authorization", "Bearer "+accessToken)
		getRecorder := httptest.NewRecorder()
		controllers.ReadAllMovies(db)(getRecorder, getReq)
		if getRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, getRecorder.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}
