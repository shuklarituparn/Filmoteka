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
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	"gorm.io/gorm"
)

func CreateMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var movie models.Movie
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			log.Error("Invalid request payload:", err.Error())
			file_logger.Println("Invalid request payload:", err.Error())
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		validate := validator.New()
		if err := validate.Struct(movie); err != nil {
			errorsMap := make(map[string]interface{})
			for _, e := range err.(validator.ValidationErrors) {
				errorsMap[e.Field()] = e.Value()
			}
			errJSON, _ := json.Marshal(errorsMap)
			common.ErrorResponse(w, http.StatusBadRequest, string(errJSON))
			return
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&movie).Error; err != nil {
				log.Error("Failed to create movie:", err.Error())
				file_logger.Println("Failed to create movie:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("%v", err.Error()))
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		w.WriteHeader(http.StatusCreated)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Movie created successfully",
			"data":    movie,
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			file_logger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func ReadAllMovies(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var totalMoviesCount int64
		var movies []models.Movie
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
			sortBy = "rating"
		}
		if sortOrder != "ASC" && sortOrder != "DESC" {
			sortOrder = "DESC"
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			query := db.Model(&movies).Limit(pageSize).Offset(offset).Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
			if err := query.Preload("Actors").Preload("Actors.Movies").Find(&movies).Error; err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
				return err
			}
			if err := tx.Model(&models.Movie{}).Select("COUNT(*)").Count(&totalMoviesCount).Error; err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		totalPages := int(math.Ceil(float64(totalMoviesCount) / float64(pageSize)))
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"data":        movies,
			"total_pages": totalPages,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func ReadMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var movie models.Movie
		movieID := r.URL.Query().Get("id")
		if movieID == "" {
			common.ErrorResponse(w, http.StatusBadRequest, "Movie ID is required")
			return
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Preload("Actors").Preload("Actors.Movies").First(&movie, movieID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					common.ErrorResponse(w, http.StatusNotFound, "Movie not found")
					return err
				}
				log.Error("Error fetching movie:", err.Error())
				file_logger.Println("Error fetching movie:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"data": movie,
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			file_logger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func UpdateMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var movie models.Movie
		var response models.Movie
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		validate := validator.New()
		if err := validate.Struct(movie); err != nil {
			errorsMap := make(map[string]interface{})
			for _, e := range err.(validator.ValidationErrors) {
				errorsMap[e.Field()] = e.Value()
			}
			errJSON, _ := json.Marshal(errorsMap)
			common.ErrorResponse(w, http.StatusBadRequest, string(errJSON))
			return
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(&movie).Error; err != nil {
				log.Error("Failed to update movie:", err.Error())
				file_logger.Println("Failed to update movie:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to update movie")
				return err
			}
			if err := tx.Preload("Actors").Preload("Actors.Movies").Find(&response, movie.ID).Error; err != nil {
				log.Error("Failed to update movie:", err.Error())
				file_logger.Println("Failed to update movie:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to update movie")
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Movie updated successfully",
			"data":    response,
		})
		if resErr != nil {
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func DeleteMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		movieID := r.URL.Query().Get("id")
		if movieID == "" {
			common.ErrorResponse(w, http.StatusBadRequest, "Movie ID is required")
			return
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("DELETE FROM actor_movies WHERE movie_id = ?", movieID).Error; err != nil {
				log.Error("Failed to delete Association:", err.Error())
				file_logger.Println("Failed to delete Association:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete Association")
				return err
			}
			if err := tx.Delete(&models.Movie{}, movieID).Error; err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete movie")
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Movie deleted successfully",
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			file_logger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func PatchMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var patchData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
			log.Error("Invalid request payload:", err.Error())
			file_logger.Println("Invalid request payload:", err.Error())
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		movieID := r.URL.Query().Get("id")
		if movieID == "" {
			common.ErrorResponse(w, http.StatusBadRequest, "Movie ID is required")
			return
		}
		var patchActors []interface{}
		if actorsInterface, ok := patchData["actors"].([]interface{}); ok {
			patchActors = actorsInterface
			delete(patchData, "actors")
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&models.Movie{}).Where("id = ?", movieID).Updates(patchData).Error; err != nil {
				log.Error("Error updating actor:", err.Error())
				file_logger.Println("Error updating actor:", err.Error())
				return err
			}
			for _, actorID := range patchActors {
				var movie models.Movie
				if err := tx.First(&movie, movieID).Error; err != nil {
					log.Error("Error finding movie:", err.Error())
					file_logger.Println("Error finding movie:", err.Error())
					return err
				}
				var actor models.Actor
				if err := tx.First(&actor, actorID).Error; err != nil {
					log.Error("Error finding actor:", err.Error())
					file_logger.Println("Error finding actor:", err.Error())
					return err
				}
				if err := tx.Model(&movie).Association("Actors").Append(&actor); err != nil {
					log.Error("Error appending actor association:", err.Error())
					file_logger.Println("Error appending actor association:", err.Error())
					return err
				}
			}
			return nil
		})
		if tx_err != nil {
			return
		}
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
