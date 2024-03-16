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
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/internal/logger"
	"github.com/shuklarituparn/Filmoteka/internal/prometheus"
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	"gorm.io/gorm"
)

var file_logger = logger.SetupLogger()

// CreateActor creates a new actor.
// @Summary Create a new actor
// @ID create-actor
// @Accept json
// @Produce json
// @Tags Actors
// @Security BearerAuth
// @Param actor body models.CreateActorModel true "Actor object to be created"
// @Success 201 {object} CreateActorResponse "Actor Added"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/actors/create [post]
func CreateActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.CreateActorApiPingCounter.Inc()
		w.Header().Set("Content-Type", "application/json")
		var actor models.Actor
		if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		if !common.ValidateAndRespond(w, actor) {
			return
		}
		for _, movie := range actor.Movies {
			if !common.ValidateAndRespond(w, movie) {
				return
			}
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

// ReadAllActor returns a list of actors with pagination support.
// @Summary Get all actors with pagination
// @ID read-all-actors
// @Produce json
// @Tags Actors
// @Param page query integer true "Page number"
// @Param page_size query integer true "Number of items per page"
// @Param sort_by query string false "Field to sort by (default birth_date)"
// @Param sort_order query string false "Sort order (ASC or DESC, default DESC)"
// @Success 200 {object} ReadAllActorResponse "List of actors"
// @Failure 400 {string} string "Invalid page_size or page"
// @Failure 500 {string} string "Internal Server Error"
// @Security BearerAuth
// @Router /api/v1/actors/all [get]
func ReadAllActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.ReadAllActorApiPingCounter.Inc()
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

// ReadActor returns details of a specific actor by ID.
// @Summary Get actor by ID
// @ID read-actor-by-id
// @Produce json
// @Tags Actors
// @Security BearerAuth
// @Param id query string true "Actor ID"
// @Success 200 {object} ReadActorResponse "Actor details"
// @Failure 400 {string} string "Actor ID is required"
// @Failure 404 {string} string "Actor not found"
// @Failure 500 {string} string "Failed to fetch actor"
// @Router /api/v1/actors/get [get]
func ReadActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.ReadOneActorApiPingCounter.Inc()
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

// UpdateActor updates an existing actor.
// @Summary Update an existing actor
// @ID update-actor
// @Accept json
// @Security BearerAuth
// @Produce json
// @Tags Actors
// @Param actor body models.UpdateActorModel true "Actor object to be updated"
// @Success 200 {object} UpdateActorResponse "Actor Updated successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Failed to update actor"
// @Router /api/v1/actors/update [put]
func UpdateActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.UpdateActorApiPingCounter.Inc()
		w.Header().Set("Content-Type", "application/json")
		var actor models.Actor
		var response models.Actor
		if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		if !common.ValidateAndRespond(w, actor) {
			return
		}
		for _, movie := range actor.Movies {
			if !common.ValidateAndRespond(w, movie) {
				return
			}
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", actor.ID).Save(&actor).Error; err != nil {
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

// DeleteActor deletes an actor and its associations from the database.
// @Summary Delete an actor
// @ID delete-actor
// @Produce json
// @Tags Actors
// @Security BearerAuth
// @Param id query string true "Actor ID"
// @Success 200 {object} DeleteActorResponse "Actor deleted successfully"
// @Failure 400 {string} string "Actor ID is required"
// @Failure 500 {string} string "Failed to delete actor or its associations"
// @Router /api/v1/actors/delete [delete]
func DeleteActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.DeleteActorApiPingCounter.Inc()
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

// PatchActor updates an existing actor with the provided patch data.
// @Summary Update an existing actor partially
// @ID patch-actor
// @Accept json
// @Security BearerAuth
// @Produce json
// @Tags Actors
// @Param id query string true "Actor ID"
// @Param patchData body models.CreateActorModel true "Patch data for updating the actor"
// @Success 200 {object} PatchActorResponse "Actor updated successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Failed to update actor or its associations"
// @Router /api/v1/actors/patch [patch]
func PatchActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.PatchActorApiPingCounter.Inc()
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
		movies, ok := patchData["movies"].([]interface{})
		if !ok {
			common.ErrorResponse(w, http.StatusBadRequest, "Movies data is missing or invalid")
			return
		}
		delete(patchData, "movies")
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			for _, movie := range movies {
				movieMap, ok := movie.(map[string]interface{})
				if !ok {
					return errors.ErrUnsupported
				}
				if err := tx.Model(&models.Movie{}).Where("id=?", movieMap["id"]).Updates(movieMap).Error; err != nil {
					log.Error("Error updating movie:", err.Error())
					file_logger.Println("Error updating movie:", err.Error())
					return err
				}
				query := "SELECT COUNT(*) FROM actor_movies WHERE actor_id = ? AND movie_id = ?"
				var count int64
				if err := tx.Raw(query, actorID, movieMap["id"]).Row().Scan(&count); err != nil {
					log.Error("Error counting actor-movie relationship:", err.Error())
					file_logger.Println("Error counting actor-movie relationship:", err.Error())
					return err
				}
				if count == 0 {
					query := "INSERT INTO actor_movies (actor_id, movie_id) VALUES (?, ?)"
					if err := tx.Exec(query, actorID, movieMap["id"]).Error; err != nil {
						log.Error("Error inserting actor-movie relationship:", err.Error())
						file_logger.Println("Error inserting actor-movie relationship:", err.Error())
						return err
					}
				}
			}
			if err := tx.Model(&models.Actor{}).Where("id = ?", actorID).Updates(patchData).Error; err != nil {
				log.Error("Error updating actor:", err.Error())
				file_logger.Println("Error updating actor:", err.Error())
				return err
			}
			return nil
		})
		if tx_err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Actor updated successfully",
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			file_logger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}
