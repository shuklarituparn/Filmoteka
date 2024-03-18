package actor_tests

import (
	"bytes"
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/pkg/hashing"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var sqliteDb *gorm.DB

func GetDBInstance() *gorm.DB {
	return sqliteDb
}

func ConnectTestDB() string {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return error.Error(err)
	}
	migrationErr := db.AutoMigrate(&models.User{}, &models.Actor{}, &models.Movie{})
	if migrationErr != nil {
		log.Error(migrationErr)
	}
	var adminUser models.User
	adminPassword, _ := hashing.HashPassword("adminpassword")
	if result := db.First(&adminUser, "role = ?", "ADMIN"); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			admin := models.User{
				Email:    "admin@example.com",
				Password: adminPassword,
				Role:     "ADMIN",
			}
			if err := db.Create(&admin).Error; err != nil {
				log.Error("Error creating admin user:", err)
			}
		} else {
			log.Error("Error querying admin user:", result.Error)
		}
	}
	sqliteDb = db
	return "Connection to test database has been established"
}

func CloseAndDelete(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic("Error getting underlying database connection: " + err.Error())
	}
	err = sqlDB.Close()
	if err != nil {
		return
	}
	err = os.Remove("test.db")
	if err != nil {
		panic("Error removing database file: " + err.Error())
	}
}

func LoginAndGetAccessToken(t *testing.T, db *gorm.DB) string {
	user := models.User{
		Email:    "admin@example.com",
		Password: "adminpassword",
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}
	loginReq, err := http.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}
	loginRecorder := httptest.NewRecorder()
	controllers.LoginUser(db)(loginRecorder, loginReq)
	if loginRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, loginRecorder.Code)
	}
	var loginResponse controllers.LoginResponse
	if err := json.NewDecoder(loginRecorder.Body).Decode(&loginResponse); err != nil {
		t.Fatal(err)
	}
	accessToken := loginResponse.AccessToken
	if accessToken == "" {
		t.Errorf("Expected access token; got empty")
	}
	return accessToken
}
