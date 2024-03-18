package actor_tests

import (
	"bytes"
	"encoding/json"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestReadAActor(t *testing.T) {
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
		t.Fatal(err)
	}
	createActor, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}
	getReq.Header.Set("Authorization", "Bearer "+accessToken)
	getRecorder := httptest.NewRecorder()
	controllers.ReadActor(db)(getRecorder, getReq)
	if getRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, getRecorder.Code)
	}
	defer CloseAndDelete(db)
}
