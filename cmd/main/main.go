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
	"github.com/shuklarituparn/Filmoteka/internal/prometheus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Filmotek API
// @version 1.0
// @description Fimotek Api Docs
// @contact.name API Support
// @contact.url https://github.com/shuklarituparn
// @contact.email rtprnshukla@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	prometheus.Init()

	err := godotenv.Load()

	var fileLogger = logger.SetupLogger()

	if err != nil {
		log.Error("Error loading .env file")
		fileLogger.Println("Error loading .env file")
	}

	config.Connect_DB()

	if os.Getenv("PORT") == "" {
		log.Error("PORT environment variable not set")
		fileLogger.Println("PORT environment variable not set")
		os.Exit(1)
	}

	ServerPort, _ := strconv.Atoi(os.Getenv("PORT"))

	port := ServerPort

	rootMux := http.NewServeMux()
	rootMux.Handle("/swagger.json", http.FileServer(http.Dir("./docs")))
	rootMux.HandleFunc("/healthcheck", controllers.HealthCheck)
	routes.ActorRouter(rootMux)
	routes.MovieRouter(rootMux)
	routes.UserRouter(rootMux)
	routes.SearchRouter(rootMux)
	rootMux.HandleFunc("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))
	log.Info(fmt.Sprintf("Server started on port %d", port))
	fileLogger.Printf("Server started on port %d", port)

	serverError := http.ListenAndServe(fmt.Sprintf(":%d", port), rootMux)
	if serverError != nil {
		log.Error("Something went wrong while starting the server!")
		fileLogger.Println("Something went wrong while starting the server!")
		panic(serverError)
	} else {
		log.Info("Server started successfully!")
		fileLogger.Println("Server started successfully!")
	}
}
