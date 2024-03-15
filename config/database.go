package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/internal/logger"
	"github.com/shuklarituparn/Filmoteka/pkg/hashing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetInstance() *gorm.DB {
	return db
}

func ConnectDb() {
	var fileLogger = logger.SetupLogger()

	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)
	portInt, _ := strconv.Atoi(port)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, portInt, user, password, dbname)

	postgresqlDb, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Error("Error connecting to the database:", err)
		fileLogger.Println("Error connecting to the database:", err)
	}
	migrationErr := postgresqlDb.AutoMigrate(&models.Actor{}, &models.User{}, &models.Movie{})
	if migrationErr != nil {
		log.Error(migrationErr)
		fileLogger.Println(migrationErr)
	}
	var adminUser models.User
	adminPassword, _ := hashing.HashPassword("adminpassword")
	if result := postgresqlDb.First(&adminUser, "role = ?", "ADMIN"); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			admin := models.User{
				Email:    "admin@example.com",
				Password: adminPassword,
				Role:     "ADMIN",
			}
			if err := postgresqlDb.Create(&admin).Error; err != nil {
				log.Error("Error creating admin user:", err)
				fileLogger.Println("Error creating admin user:", err)
			}
		} else {
			log.Error("Error querying admin user:", result.Error)
			fileLogger.Println("Error querying admin user:", result.Error)
		}
	}
	db = postgresqlDb
	log.Info("Successfully connected!")
	fileLogger.Println("Successfully connected to the database!")
}
