package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

// HealthCheck performs a health check and returns the status of the application.
// @Summary Perform health check
// @Tags Healthcheck
// @ID health-check
// @Produce json
// @Success 200 {object} HealthCheckResponse "Health check response"
// @Router /healthcheck [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Author:      "Rituparn Shukla",
		CurrentTime: time.Now(),
		Status:      "up",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error encoding JSON:", err.Error())
		fileLogger.Println("Error encoding JSON:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
