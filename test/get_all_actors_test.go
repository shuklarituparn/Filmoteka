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

func TestReadAllActors(t *testing.T) {

	r := runner.NewRunner(t, "Actor")
	r.NewTest("Actor Reading", func(P provider.T) {
		P.Epic("Actor controller")
		P.Feature("Actor getting")
		P.Title("Getting an actor")
		P.Description("This test checks if the actor can be read or not")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		getReq, err := http.NewRequest("GET", "/api/v1/actors/all?page_size=3&page=1", bytes.NewReader([]byte{}))
		if err != nil {
			t.Error(err)
		}
		getReq.Header.Set("Authorization", "Bearer "+accessToken)
		getRecorder := httptest.NewRecorder()
		controllers.ReadAllActor(db)(getRecorder, getReq)
		if getRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, getRecorder.Code)
		}
		defer CloseAndDelete(db)
	})

	r.RunTests()

}

func TestReadAllActorsWithoutPagination(t *testing.T) {
	r := runner.NewRunner(t, "Actor")
	r.NewTest("Actor Reading", func(P provider.T) {
		P.Epic("Actor controller")
		P.Feature("Actor getting")
		P.Title("Getting all actor")
		P.Description("This test checks if the all the actors can be read or not without pagination")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		getReq, err := http.NewRequest("GET", "/api/v1/actors/all", bytes.NewReader([]byte{}))
		if err != nil {
			t.Error(err)
		}
		getReq.Header.Set("Authorization", "Bearer "+accessToken)
		getRecorder := httptest.NewRecorder()
		controllers.ReadAllActor(db)(getRecorder, getReq)
		if getRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, getRecorder.Code)
		}
		defer CloseAndDelete(db)
	})

	r.RunTests()

}
