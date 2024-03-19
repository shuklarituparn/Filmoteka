package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateActor(t *testing.T) {
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
	createActorReq, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
	if err != nil {
		t.Error(err)
	}
	createActorReq.Header.Set("Authorization", "Bearer "+accessToken)
	createActorRecorder := httptest.NewRecorder()
	CreateActor(db)(createActorRecorder, createActorReq)
	if createActorRecorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, createActorRecorder.Code)
	}
	var createActorResponse CreateActorResponse
	if err := json.NewDecoder(createActorRecorder.Body).Decode(&createActorResponse); err != nil {
		t.Error(err)
	}
	updateData := models.Actor{
		ID:        uint(createActorResponse.ID),
		LastName:  "UpdatedLastName",
		FirstName: "UpdatedFirstName",
		Gender:    "Male",
		BirthDate: "2000-10-10",
	}
	updateActorJSON, err := json.Marshal(updateData)
	if err != nil {
		t.Error(err)
	}
	updateActorReq, err := http.NewRequest("PUT", "/api/v1/actors/update", bytes.NewBuffer(updateActorJSON))
	if err != nil {
		t.Error(err)
	}
	updateActorReq.Header.Set("Authorization", "Bearer "+accessToken)
	updateActorRecorder := httptest.NewRecorder()
	UpdateActor(db)(updateActorRecorder, updateActorReq)
	if updateActorRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, updateActorRecorder.Code)
	}
	defer CloseAndDelete(db)
}
