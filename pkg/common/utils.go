package common

import (
	"encoding/json"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/shuklarituparn/Filmoteka/internal/logger"
)

var fileLogger = logger.SetupLogger()

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	log.Error(message)
	fileLogger.Println(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errorMsg := map[string]string{"error": message}
	err := json.NewEncoder(w).Encode(errorMsg)
	if err != nil {
		log.Error("Error encoding JSON:", err.Error())
		fileLogger.Println("Error encoding JSON:", err.Error())
	}
}
