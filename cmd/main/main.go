package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/api/routes"
	"github.com/shuklarituparn/Filmoteka/config"
	"github.com/shuklarituparn/Filmoteka/internal/logger"
)

func main() {
	err := godotenv.Load()

	var file_logger = logger.SetupLogger()

	if err != nil {
		log.Error("Error loading .env file")
		file_logger.Println("Error loading .env file")
	}

	config.ConnectDb()

	if os.Getenv("PORT") == "" {
		log.Error("PORT environment variable not set")
		file_logger.Println("PORT environment variable not set")
		os.Exit(1)
	}

	ServerPort, _ := strconv.Atoi(os.Getenv("PORT"))

	port := ServerPort

	rootMux := http.NewServeMux()

	rootMux.HandleFunc("/healthcheck", controllers.HealthCheck)
	routes.ActorRouter(rootMux)
	routes.MovieRouter(rootMux)
	routes.UserRouter(rootMux)
	routes.SearchRouter(rootMux)

	log.Info(fmt.Sprintf("Server started on port %d", port))
	file_logger.Printf("Server started on port %d", port)

	serverError := http.ListenAndServe(fmt.Sprintf(":%d", port), rootMux)
	if serverError != nil {
		log.Error("Something went wrong while starting the server!")
		file_logger.Println("Something went wrong while starting the server!")
		panic(serverError)
	} else {
		log.Info("Server started successfully!")
		file_logger.Println("Server started successfully!")
	}
}
