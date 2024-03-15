package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

type HealthCheckResponse struct {
	Author      string    `json:"author"`
	CurrentTime time.Time `json:"current_time"`
	Status      string    `json:"status"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Author:      "Rituparn Shukla",
		CurrentTime: time.Now(),
		Status:      "up",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Error encoding JSON:", err.Error())
		file_logger.Println("Error encoding JSON:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
