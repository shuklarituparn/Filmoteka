package user_tests

import (
	"bytes"
	"encoding/json"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"net/http"
	"net/http/httptest"
	"testing"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestLogin(t *testing.T) {
	ConnectTestDB()
	var db = GetDBInstance()
	user := User{
		Email:    "admin@example.com",
		Password: "adminpassword",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	controllers.LoginUser(db)(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}
	defer CloseAndDelete(db)
}
