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

func TestPatchActor(t *testing.T) {
	r := runner.NewRunner(t, "Actor")
	r.NewTest("Checking Actor Patch", func(P provider.T) {
		P.Feature("Actor Patch")
		P.Title("Testing actor patch")
		P.Description("This test checks if the actor can be patched or not")
		ConnectTestDB()
		var db = GetDBInstance()
		defer CloseAndDelete(db)
		accessToken := LoginAndGetAccessToken(t, db)
		actor := models.Actor{
			BirthDate: "1998-05-28",
			FirstName: "Rituparn",
			Gender:    "Male",
			LastName:  "Shukla",
		}
		actorJSON, err := json.Marshal(actor)
		if err != nil {
			t.Error(err)
		}
		createActorReq, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
		if err != nil {
			t.Error(err)
		}
		createActorReq.Header.Set("Authorization", "Bearer "+accessToken)
		createActorRecorder := httptest.NewRecorder()
		controllers.CreateActor(db)(createActorRecorder, createActorReq)
		if createActorRecorder.Code != http.StatusCreated {
			t.Errorf("Expected status %d; got %d", http.StatusCreated, createActorRecorder.Code)
		}
		var createActorResponse controllers.CreateActorResponse
		if err := json.NewDecoder(createActorRecorder.Body).Decode(&createActorResponse); err != nil {
			t.Error(err)
		}
		patchData := models.Actor{
			LastName:  "PatchedLastName",
			FirstName: "PatchedFirstName",
		}
		patchActorJSON, err := json.Marshal(patchData)
		if err != nil {
			t.Error(err)
		}
		patchActorReq, err := http.NewRequest("PATCH", "/api/v1/actors/patch?id="+strconv.Itoa(createActorResponse.ID), bytes.NewBuffer(patchActorJSON))
		if err != nil {
			t.Error(err)
		}
		patchActorReq.Header.Set("Authorization", "Bearer "+accessToken)
		patchActorRecorder := httptest.NewRecorder()
		controllers.PatchActor(db)(patchActorRecorder, patchActorReq)
		fmt.Print(patchActorRecorder.Body)
		if patchActorRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, patchActorRecorder.Code)
		}

	})
	r.RunTests()

}

func TestPatchActorWithoutBody(t *testing.T) {
	r := runner.NewRunner(t, "Actor")
	r.NewTest("Checking Actor Patch", func(P provider.T) {
		P.Feature("Actor Patch")
		P.Title("Testing actor patch without body")
		P.Description("This test checks if the actor can be patched or not without body")
		ConnectTestDB()
		var db = GetDBInstance()
		defer CloseAndDelete(db)
		accessToken := LoginAndGetAccessToken(t, db)
		actor := models.CreateActorModel{
			BirthDate: "1998-05-28",
			FirstName: "Rituparn",
			Gender:    "Male",
			LastName:  "Shukla",
		}
		actorJSON, err := json.Marshal(actor)
		if err != nil {
			t.Error(err)
		}
		createActorReq, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
		if err != nil {
			t.Error(err)
		}
		createActorReq.Header.Set("Authorization", "Bearer "+accessToken)
		createActorRecorder := httptest.NewRecorder()
		controllers.CreateActor(db)(createActorRecorder, createActorReq)
		if createActorRecorder.Code != http.StatusCreated {
			t.Errorf("Expected status %d; got %d", http.StatusCreated, createActorRecorder.Code)
		}
		var createActorResponse controllers.CreateActorResponse
		if err := json.NewDecoder(createActorRecorder.Body).Decode(&createActorResponse); err != nil {
			t.Error(err)
		}
		patchData := models.CreateMovieModel{}
		patchActorJSON, err := json.Marshal(patchData)
		if err != nil {
			t.Error(err)
		}
		patchActorReq, err := http.NewRequest("PATCH", "/api/v1/actors/patch?id="+strconv.Itoa(createActorResponse.ID), bytes.NewBuffer(patchActorJSON))
		if err != nil {
			t.Error(err)
		}
		patchActorReq.Header.Set("Authorization", "Bearer "+accessToken)
		patchActorRecorder := httptest.NewRecorder()
		controllers.PatchActor(db)(patchActorRecorder, patchActorReq)
		fmt.Print(patchActorRecorder.Body)
		if patchActorRecorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, patchActorRecorder.Code)
		}

	})
	r.RunTests()

}
