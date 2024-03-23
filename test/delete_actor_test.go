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

func TestDeleteActor(t *testing.T) {

	r := runner.NewRunner(t, "Actor")
	r.NewTest("Checking actor deletion", func(P provider.T) {
		P.Feature("Actor Deletion")
		P.Title("Testing actor deletion")
		P.Description("This test checks if the actor can be deleted")
		ConnectTestDB()
		var db = GetDBInstance()
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
		fmt.Println(actorResponse.ID)
		deleteReq, err := http.NewRequest("DELETE", "/api/v1/actors/delete?id="+strconv.Itoa(actorResponse.ID), bytes.NewReader([]byte{}))
		if err != nil {
			t.Error(err)
		}
		deleteReq.Header.Set("Authorization", "Bearer "+accessToken)
		deleteRecorder := httptest.NewRecorder()
		controllers.DeleteActor(db)(deleteRecorder, deleteReq)
		if deleteRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, deleteRecorder.Code)
		}
		defer CloseAndDelete(db)

	})
	r.RunTests()

}

func TestDeleteActorWithoutBody(t *testing.T) {

	r := runner.NewRunner(t, "Actor")
	r.NewTest("Checking actor deletion", func(P provider.T) {
		P.Feature("Actor Deletion")
		P.Title("Testing actor deletion without a body")
		P.Description("This test checks if the actor can be deleted without a body or not")
		ConnectTestDB()
		var db = GetDBInstance()
		accessToken := LoginAndGetAccessToken(t, db)
		actor := models.CreateActorModel{}
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
		fmt.Println(actorResponse.ID)
		deleteReq, err := http.NewRequest("DELETE", "/api/v1/actors/delete?id="+strconv.Itoa(actorResponse.ID), bytes.NewReader([]byte{}))
		if err != nil {
			t.Error(err)
		}
		deleteReq.Header.Set("Authorization", "Bearer "+accessToken)
		deleteRecorder := httptest.NewRecorder()
		controllers.DeleteActor(db)(deleteRecorder, deleteReq)
		if deleteRecorder.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, deleteRecorder.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}
