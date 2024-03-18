package movie_tests

import (
	"bytes"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadAllMovies(t *testing.T) {
	ConnectTestDB()
	var db = GetDBInstance()
	accessToken := LoginAndGetAccessToken(t, db)
	getReq, err := http.NewRequest("GET", "/api/v1/movies/all?page_size=3&page=1", bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("Authorization", "Bearer "+accessToken)
	getRecorder := httptest.NewRecorder()
	controllers.ReadAllMovies(db)(getRecorder, getReq)
	if getRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, getRecorder.Code)
	}
	defer CloseAndDelete(db)
}