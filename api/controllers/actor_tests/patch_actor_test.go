package actor_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestPatchActor(t *testing.T) {
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
		t.Fatal(err)
	}
	createActorReq, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
	if err != nil {
		t.Fatal(err)
	}
	createActorReq.Header.Set("Authorization", "Bearer "+accessToken)
	createActorRecorder := httptest.NewRecorder()
	controllers.CreateActor(db)(createActorRecorder, createActorReq)
	if createActorRecorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, createActorRecorder.Code)
	}
	var createActorResponse controllers.CreateActorResponse
	if err := json.NewDecoder(createActorRecorder.Body).Decode(&createActorResponse); err != nil {
		t.Fatal(err)
	}
	patchData := models.Actor{
		LastName:  "PatchedLastName",
		FirstName: "PatchedFirstName",
	}
	patchActorJSON, err := json.Marshal(patchData)
	if err != nil {
		t.Fatal(err)
	}
	patchActorReq, err := http.NewRequest("PATCH", "/api/v1/actors/patch?id="+strconv.Itoa(createActorResponse.ID), bytes.NewBuffer(patchActorJSON))
	if err != nil {
		t.Fatal(err)
	}
	patchActorReq.Header.Set("Authorization", "Bearer "+accessToken)
	patchActorRecorder := httptest.NewRecorder()
	controllers.PatchActor(db)(patchActorRecorder, patchActorReq)
	fmt.Print(patchActorRecorder.Body)
	if patchActorRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, patchActorRecorder.Code)
	}
}
