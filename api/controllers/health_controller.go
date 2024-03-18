package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

// HealthCheck выполняет проверку состояния и возвращает статус приложения.
// @Summary Выполнить проверку состояния
// @Tags Healthcheck
// @ID health-check
// @Produce json
// @Success 200 {object} HealthCheckResponse "Ответ на проверку состояния"
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
