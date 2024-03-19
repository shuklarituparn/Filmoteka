package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	ConnectTestDB()
	var db = GetDBInstance()
	user := User{
		Email:    "rituparn@example.com",
		Password: "rituparn",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	RegisterUser(db)(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}
	defer CloseAndDelete(db)
}
