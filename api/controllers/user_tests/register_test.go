package user_tests

import (
	"bytes"
	"encoding/json"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
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
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	controllers.RegisterUser(db)(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}
	defer CloseAndDelete(db)
}
