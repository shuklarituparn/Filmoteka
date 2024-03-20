package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"io"
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
		t.Error(err)
	}
	createReq, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(actorJSON))
	if err != nil {
		t.Error(err)
	}
	createReq.Header.Set("Authorization", "Bearer "+accessToken)
	createRecorder := httptest.NewRecorder()
	CreateActor(db)(createRecorder, createReq)
	if createRecorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusOK, createRecorder.Code)
	}
	defer CloseAndDelete(db)
}

func TestCreateActor_Failure(t *testing.T) {
	// Simulate a failure scenario for invalid request payload
	t.Run("Invalid_Request_Payload", func(t *testing.T) {
		// Create a request with invalid JSON payload
		invalidPayload := []byte("invalid JSON")
		req, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(invalidPayload))
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		CreateActor(nil)(recorder, req)

		// Check if the response status is BadRequest (400)
		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	// Simulate a failure scenario for internal server error while closing request body
	t.Run("Internal_Server_Error_On_Body_Close", func(t *testing.T) {
		// Create a request with valid JSON payload
		validPayload := []byte(`{"firstName":"John","lastName":"Doe","gender":"Male","birthDate":"1990-01-01","movies":[]}`)
		req, err := http.NewRequest("POST", "/api/v1/actors/create", bytes.NewBuffer(validPayload))
		if err != nil {
			t.Fatal(err)
		}

		// Mock the Body to simulate an error while closing
		req.Body = &mockFailingBody{}

		recorder := httptest.NewRecorder()
		CreateActor(nil)(recorder, req)

		// Check if the response status is InternalServerError (500)
		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status %d; got %d", http.StatusInternalServerError, recorder.Code)
		}
	})

	// Add more failure scenarios as needed for remaining failure points
}

// Mocking a failing ReadCloser to simulate an error while closing the request body
type mockFailingBody struct{}

func (*mockFailingBody) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (*mockFailingBody) Close() error {
	return errors.New("forced error")
}
