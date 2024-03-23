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

func TestRegister(t *testing.T) {

	r := runner.NewRunner(t, "User")
	r.NewTest("User Register", func(P provider.T) {
		P.Feature("Checking user registration")
		P.Title("Testing if user can register or not")
		P.Description("This test checks if the users can register or not")
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
		controllers.RegisterUser(db)(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}

func TestRegisterWithoutBody(t *testing.T) {

	r := runner.NewRunner(t, "User")
	r.NewTest("User Register", func(P provider.T) {
		P.Feature("Checking user registration")
		P.Title("Testing if user can register or not without body")
		P.Description("This test checks if the users can register or not without body")
		ConnectTestDB()
		var db = GetDBInstance()
		user := User{}
		userJSON, err := json.Marshal(user)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Error(err)
		}
		rr := httptest.NewRecorder()
		controllers.RegisterUser(db)(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d; got %d", http.StatusBadRequest, rr.Code)
		}
		defer CloseAndDelete(db)
	})
	r.RunTests()

}
