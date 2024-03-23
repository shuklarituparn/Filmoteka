package test

import (
	"bytes"
	"encoding/json"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
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

	r := runner.NewRunner(t, "User")
	r.NewTest("User login", func(P provider.T) {
		P.Feature("Checking user login")
		P.Title("Testing if user can login or not")
		P.Description("This test checks if the users can login or not")
		ConnectTestDB()
		var db = GetDBInstance()
		user := User{
			Email:    "admin@example.com",
			Password: "adminpassword",
		}
		userJSON, err := json.Marshal(user)
		if err != nil {
			t.Error(err)
		}
		req, err := http.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		controllers.LoginUser(db)(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}

func TestLoginWithoutBody(t *testing.T) {
	r := runner.NewRunner(t, "User")
	r.NewTest("User login", func(P provider.T) {
		P.Feature("Checking user login")
		P.Title("Testing if user can login or not without body")
		P.Description("This test checks if the users can login or not without body")
		ConnectTestDB()
		var db = GetDBInstance()
		user := User{}
		userJSON, err := json.Marshal(user)
		if err != nil {
			t.Error(err)
		}
		req, err := http.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		controllers.LoginUser(db)(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d; got %d", http.StatusUnauthorized, rr.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}
