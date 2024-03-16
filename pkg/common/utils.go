package common

import (
	"encoding/json"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/internal/logger"
)

var file_logger = logger.SetupLogger()

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	log.Error(message)
	file_logger.Println(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errorMsg := map[string]string{"error": message}
	err := json.NewEncoder(w).Encode(errorMsg)
	if err != nil {
		log.Error("Error encoding JSON:", err.Error())
		file_logger.Println("Error encoding JSON:", err.Error())
	}
}

func ValidateAndRespond(w http.ResponseWriter, v interface{}) bool {
	validate := validator.New()
	err := validate.RegisterValidation("validDate", models.ValidDate)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Server Error!")
	}
	if err := validate.Struct(v); err != nil {
		errorsMap := make(map[string]interface{})
		for _, e := range err.(validator.ValidationErrors) {
			errorsMap[e.Field()] = e.Tag()
		}
		errJSON, _ := json.Marshal(errorsMap)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(errJSON)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Server Error!")
		}
		return false
	}
	return true
}
