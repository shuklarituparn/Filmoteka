package actor_tests

import (
	"bytes"
	"encoding/json"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateActor(t *testing.T) {
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
	createReq, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
	if err != nil {
		t.Fatal(err)
	}
	createReq.Header.Set("Authorization", "Bearer "+accessToken)
	createRecorder := httptest.NewRecorder()
	controllers.CreateActor(db)(createRecorder, createReq)
	if createRecorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusOK, createRecorder.Code)
	}
	defer CloseAndDelete(db)
}
