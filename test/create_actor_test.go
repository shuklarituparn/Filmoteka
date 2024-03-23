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

func TestCreateActor(t *testing.T) {

	r := runner.NewRunner(t, "Actor")
	r.NewTest("Actor creation", func(P provider.T) {
		P.Epic("Actor controller")
		P.Feature("Actor Creation")
		P.Title("Testing actor creation")
		P.Description("This test checks if the actor can be created or not")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		actor := models.CreateActorModel{
			BirthDate: "1998-05-28",
			FirstName: "Rituparn",
			Gender:    "Male",
			LastName:  "Shukla",
		}
		actorJSON, err := json.Marshal(actor)
		if err != nil {
			P.Error(err)
		}
		createReq, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
		if err != nil {
			P.Error(err)
		}
		createReq.Header.Set("Authorization", "Bearer "+accessToken)
		createRecorder := httptest.NewRecorder()
		controllers.CreateActor(db)(createRecorder, createReq)
		if createRecorder.Code != http.StatusCreated {
			P.Errorf("Expected status %d; got %d", http.StatusOK, createRecorder.Code)
		}
		defer CloseAndDelete(db)
	})

	r.RunTests()

}

func TestCreateActorWithoutBody(t *testing.T) {

	r := runner.NewRunner(t, "Actor")
	r.NewTest("Checking false actor creation", func(P provider.T) {
		P.Feature("Actor Creation")
		P.Title("Testing actor creation")
		P.Description("This test checks if the actor can be created or not without a body")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		actor := models.Actor{
			BirthDate: "2003-10-10",
		}
		actorJSON, err := json.Marshal(actor)
		if err != nil {
			t.Error(err)
		}
		createReq, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
		if err != nil {
			t.Error(err)
		}
		createReq.Header.Set("Authorization", "Bearer "+accessToken)
		createRecorder := httptest.NewRecorder()
		controllers.CreateActor(db)(createRecorder, createReq)
		if createRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, createRecorder.Code)
		}
		defer CloseAndDelete(db)

	})
	r.RunTests()
}
