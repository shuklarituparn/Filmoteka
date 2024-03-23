package test

import (
	"encoding/json"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	r := runner.NewRunner(t, "Healthcheck")
	r.NewTest("System status", func(P provider.T) {
		P.Feature("Checking system status")
		P.Title("Testing if the system is up")
		P.Description("This test checks if the service is up or not")
		req, err := http.NewRequest("GET", "/healthcheck", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(controllers.HealthCheck)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		var response controllers.HealthCheckResponse
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Error(err)
		}
		if response.Author != "Rituparn Shukla" {
			t.Errorf("Expected author %s; got %s", "Rituparn Shukla", response.Author)
		}
		if response.Status != "up" {
			t.Errorf("Expected status %s; got %s", "up", response.Status)
		}
	})
	r.RunTests()

}
