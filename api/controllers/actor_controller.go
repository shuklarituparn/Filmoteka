package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/internal/logger"
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	"gorm.io/gorm"
)

var file_logger = logger.SetupLogger()

func CreateActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var actor models.Actor
		if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		validate := validator.New()
		if err := validate.Struct(actor); err != nil {
			errorsMap := make(map[string]interface{})
			for _, e := range err.(validator.ValidationErrors) {
				errorsMap[e.Field()] = e.Value()
			}
			errJSON, _ := json.Marshal(errorsMap)
			common.ErrorResponse(w, http.StatusBadRequest, string(errJSON))
			return
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&actor).Error; err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to create actor")
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		w.WriteHeader(http.StatusCreated)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{"id": actor.ID, "message": "Actor Added"})
		if resErr != nil {
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func ReadAllActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var actors []models.Actor
		var totalActorsCount int64
		pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
		if err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid page_size")
			return
		}
		pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || pageNum < 1 {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid page")
			return
		}
		offset := (pageNum - 1) * pageSize
		sortBy := r.URL.Query().Get("sort_by")
		sortOrder := strings.ToUpper(r.URL.Query().Get("sort_order"))
		if sortBy == "" {
			sortBy = "birth_date"
		}
		if sortOrder != "ASC" && sortOrder != "DESC" {
			sortOrder = "DESC"
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&models.Actor{}).Select("COUNT(*)").Count(&totalActorsCount).Error; err != nil {
				log.Error("Error counting actors:", err.Error())
				file_logger.Println("Error counting actors:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
				return err
			}
			query := tx.Model(&actors).Limit(pageSize).Offset(offset).Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
			if err := query.Preload("Movies").Preload("Movies.Actors").Find(&actors).Error; err != nil {
				log.Error("Error fetching actors:", err.Error())
				file_logger.Println("Error fetching actors:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		totalPages := int(math.Ceil(float64(totalActorsCount) / float64(pageSize)))
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": actors, "total_pages": totalPages}); err != nil {
			log.Error("Error encoding JSON:", err.Error())
			file_logger.Println("Error encoding JSON:", err.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func ReadActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var actor models.Actor
		actorID := r.URL.Query().Get("id")
		if actorID == "" {
			common.ErrorResponse(w, http.StatusBadRequest, "Actor ID is required")
			return
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Preload("Movies").Preload("Movies.Actors").First(&actor, actorID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					common.ErrorResponse(w, http.StatusNotFound, "Actor not found")
					return err
				}
				log.Error("Error fetching actor:", err.Error())
				file_logger.Println("Error fetching actor:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch actor")
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{"data": actor})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			file_logger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func UpdateActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var actor models.Actor
		var response models.Actor
		if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		validate := validator.New()
		if err := validate.Struct(actor); err != nil {
			errorsMap := make(map[string]interface{})
			for _, e := range err.(validator.ValidationErrors) {
				errorsMap[e.Field()] = e.Value()
			}
			errJSON, _ := json.Marshal(errorsMap)
			common.ErrorResponse(w, http.StatusBadRequest, string(errJSON))
			return
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(&actor).Error; err != nil {
				log.Error("Error updating actor:", err.Error())
				file_logger.Println("Error updating actor:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to update actor")
				return err
			}
			if err := tx.Preload("Movies").Preload("Movies.Actors").Find(&response, actor.ID).Error; err != nil {
				log.Error("Error updating actor:", err.Error())
				file_logger.Println("Error updating actor:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to update actor")
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{"actor": response, "message": "Actor Updated successfully"})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			file_logger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func DeleteActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		actorID := r.URL.Query().Get("id")
		if actorID == "" {
			common.ErrorResponse(w, http.StatusBadRequest, "Actor ID is required")
			return
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("DELETE FROM actor_movies WHERE actor_id = ?", actorID).Error; err != nil {
				log.Error("Error deleting association:", err.Error())
				file_logger.Println("Error deleting association:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete Association")
				return err
			}
			if err := tx.Delete(&models.Actor{}, actorID).Error; err != nil {
				log.Error("Error deleting actor:", err.Error())
				file_logger.Println("Error deleting actor:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete actor")
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Actor deleted successfully",
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			file_logger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func PatchActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var patchData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		actorID := r.URL.Query().Get("id")
		if actorID == "" {
			common.ErrorResponse(w, http.StatusBadRequest, "Actor ID is required")
			return
		}
		var patchMovies []interface{}
		if moviesInterface, ok := patchData["movies"].([]interface{}); ok {
			patchMovies = moviesInterface
			delete(patchData, "movies")
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&models.Actor{}).Where("id = ?", actorID).Updates(patchData).Error; err != nil {
				log.Error("Error updating actor:", err.Error())
				file_logger.Println("Error updating actor:", err.Error())
				return err
			}
			for _, movieID := range patchMovies {
				var actor models.Actor
				if err := tx.First(&actor, actorID).Error; err != nil {
					log.Error("Error finding actor:", err.Error())
					file_logger.Println("Error finding actor:", err.Error())
					return err
				}
				var movie models.Movie
				if err := tx.First(&movie, movieID).Error; err != nil {
					log.Error("Error finding movie:", err.Error())
					file_logger.Println("Error finding movie:", err.Error())
					return err
				}
				if err := tx.Model(&actor).Association("Movies").Append(&movie); err != nil {
					log.Error("Error appending movie association:", err.Error())
					file_logger.Println("Error appending movie association:", err.Error())
					return err
				}
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Movie updated successfully",
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			file_logger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}
