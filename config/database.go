package config

import (
	"fmt"
	"github.com/pkg/errors"
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
	var file_logger = logger.SetupLogger()

	var (
		host     = os.Getenv("POSTGRES_HOST")
		port     = os.Getenv("POSTGRES_PORT")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
	)
	portInt, _ := strconv.Atoi(port)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, portInt, user, password, dbname)

	postgresqlDb, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Error("Error connecting to the database:", err)
		file_logger.Println("Error connecting to the database:", err)
	}
	migrationErr := postgresqlDb.AutoMigrate(&models.Actor{}, &models.User{}, &models.Movie{})
	if migrationErr != nil {
		log.Error(migrationErr)
		file_logger.Println(migrationErr)
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
				file_logger.Println("Error creating admin user:", err)
			}
		} else {
			log.Error("Error querying admin user:", result.Error)
			file_logger.Println("Error querying admin user:", result.Error)
		}
	}
	db = postgresqlDb
	log.Info("Successfully connected!")
	file_logger.Println("Successfully connected to the database!")
}
