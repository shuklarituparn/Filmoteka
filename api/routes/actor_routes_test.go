package routes

//import (
//	"fmt"
//	"github.com/charmbracelet/log"
//	jwt "github.com/shuklarituparn/Filmoteka/pkg/jwt_token"
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//	"time"
//)
//
//var getActorRoute = fmt.Sprintf("/api/v1/actors/get?id=%d", 3)
//
//func TestActorRouter(t *testing.T) {
//	mux := http.NewServeMux()
//	ActorRouter(mux)
//
//	psqlInfo := fmt.Sprintf("host=localhost port=5432 user=rituparn " +
//		"password=rituparn28 dbname=Filmotek sslmode=disable")
//
//	postgresqlDb, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{}) //we get *db.db
//	if err != nil {
//		log.Error("Error opening the database for connection")
//	}
//	testCases := []struct {
//		name       string
//		method     string
//		route      string
//		statusCode int
//	}{
//		{
//			name:       "GetActorHandler",
//			method:     http.MethodGet,
//			route:      getActorRoute,
//			statusCode: http.StatusOK,
//		},
//		{
//			name:       "AllActorsHandler",
//			method:     http.MethodGet,
//			route:      "/api/v1/actors/all?page=1&page_size=1",
//			statusCode: http.StatusOK,
//		},
//		// Add more test cases for other routes as needed
//	}
//	Token, _ := jwt.GetJWTToken("admin@example.com", "ADMIN", time.Hour)
//	TokenString := fmt.Sprintf("Bearer %s", Token)
//	// Iterate over test cases
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			// Create a new HTTP request with the specified method and route
//			req, err := http.NewRequest(tc.method, tc.route, nil)
//			req.Header.Set("Authorization", TokenString) //so the token is going correctly
//			if err != nil {
//				t.Fatalf("could not create request: %v", err)
//			}
//
//			// Create a new recorder to capture the response
//			rec := httptest.NewRecorder()
//
//			mux.ServeHTTP(rec, req)
//
//			// Check if the response status code matches the expected status code
//			if rec.Code != tc.statusCode {
//				t.Errorf("expected status code %d, got %d", tc.statusCode, rec.Code)
//			}
//
//			// Example assertion: Check if the response body contains certain string
//			if !strings.Contains(rec.Body.String(), "Expected response body content") {
//				t.Errorf("expected response body to contain certain content")
//			}
//		})
//	}
//}

//NEED TO GENERATE THE TOKEN WITH CLAIMS AND THEN CHECK THE CLAIMS, IF I GENERATE HERE,
//I guess the problem is that when I call the handlers the DB is not initialized
