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

func TestReadAActor(t *testing.T) {

	r := runner.NewRunner(t, "Actor")
	r.NewTest("Actor Reading", func(P provider.T) {
		P.Epic("Actor controller")
		P.Feature("Actor getting")
		P.Title("Getting an actor")
		P.Description("This test checks if the actor can be read or not")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		actor := models.Actor{
			BirthDate: "2100-05-28",
			FirstName: "Rituparn",
			Gender:    "Male",
			LastName:  "Shukla",
		}
		actorJSON, err := json.Marshal(actor)
		if err != nil {
			t.Error(err)
		}
		createActor, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
		if err != nil {
			t.Error(err)
		}
		createActor.Header.Set("Authorization", "Bearer "+accessToken)
		createRecorder := httptest.NewRecorder()
		controllers.CreateActor(db)(createRecorder, createActor)
		var actorResponse controllers.CreateActorResponse
		err = json.NewDecoder(createRecorder.Body).Decode(&actorResponse)
		if err != nil {
			return
		}
		getReq, err := http.NewRequest("GET", "/api/v1/actors/get?id="+strconv.Itoa(actorResponse.ID), bytes.NewReader([]byte{}))
		if err != nil {
			t.Error(err)
		}
		getReq.Header.Set("Authorization", "Bearer "+accessToken)
		getRecorder := httptest.NewRecorder()
		controllers.ReadActor(db)(getRecorder, getReq)
		if getRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, getRecorder.Code)
		}
		defer CloseAndDelete(db)
	})

	r.RunTests()

}

func TestReadAActorWithWrongID(t *testing.T) {

	r := runner.NewRunner(t, "Actor")
	r.NewTest("Actor Reading", func(P provider.T) {
		P.Epic("Actor controller")
		P.Feature("Actor getting")
		P.Title("Getting an actor")
		P.Description("This test checks if the actor can be read or not without ID")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		actor := models.Actor{
			BirthDate: "2100-05-28",
			FirstName: "Rituparn",
			Gender:    "Male",
			LastName:  "Shukla",
		}
		actorJSON, err := json.Marshal(actor)
		if err != nil {
			t.Error(err)
		}
		createActor, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
		if err != nil {
			t.Error(err)
		}
		createActor.Header.Set("Authorization", "Bearer "+accessToken)
		createRecorder := httptest.NewRecorder()
		controllers.CreateActor(db)(createRecorder, createActor)
		var actorResponse controllers.CreateActorResponse
		err = json.NewDecoder(createRecorder.Body).Decode(&actorResponse)
		if err != nil {
			return
		}
		getReq, err := http.NewRequest("GET", "/api/v1/actors/get?id=", bytes.NewReader([]byte{}))
		if err != nil {
			t.Error(err)
		}
		getReq.Header.Set("Authorization", "Bearer "+accessToken)
		getRecorder := httptest.NewRecorder()
		controllers.ReadActor(db)(getRecorder, getReq)
		if getRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, getRecorder.Code)
		}
		defer CloseAndDelete(db)
	})

	r.RunTests()

}
