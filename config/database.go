package config

import (
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

func Get_Instance() *gorm.DB {
	return db
}

func Connect_DB() {
	var file_logger = logger.SetupLogger()

	var (
		host     = os.Getenv("POSTGRES_HOST")
		port     = os.Getenv("POSTGRES_PORT")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
	)
	port_int, _ := strconv.Atoi(port)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port_int, user, password, dbname)

	postgresql_db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Error("Error connecting to the database:", err)
		file_logger.Println("Error connecting to the database:", err)
	}
	migration_err := postgresql_db.AutoMigrate(&models.Actor{}, &models.User{}, &models.Movie{})
	if migration_err != nil {
		log.Error(migration_err)
		file_logger.Println(migration_err)
	}
	var adminUser models.User
	admin_password, _ := hashing.HashPassword("adminpassword")
	if result := postgresql_db.First(&adminUser, "role = ?", "ADMIN"); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			admin := models.User{
				Email:    "admin@example.com",
				Password: admin_password,
				Role:     "ADMIN",
			}
			if err := postgresql_db.Create(&admin).Error; err != nil {
				log.Error("Error creating admin user:", err)
				file_logger.Println("Error creating admin user:", err)
			}
		} else {
			log.Error("Error querying admin user:", result.Error)
			file_logger.Println("Error querying admin user:", result.Error)
		}
	}
	db = postgresql_db
	log.Info("Successfully connected!")
	file_logger.Println("Successfully connected to the database!")
}
