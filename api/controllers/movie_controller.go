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
	"github.com/shuklarituparn/Filmoteka/internal/prometheus"
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	"gorm.io/gorm"
)

// CreateMovie creates a new movie.
// @Summary Create a new movie
// @ID create-movie
// @Security BearerAuth
// @Accept json
// @Produce json
// @Tags Movies
// @Param movie body models.CreateMovieModel true "Movie object to be created"
// @Success 201 {object} CreateMovieResponse "Movie created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/movies/create [post]
func CreateMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.CreateMovieApiPingCounter.Inc()
		w.Header().Set("Content-Type", "application/json")
		var movie models.Movie
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			log.Error("Invalid request payload:", err.Error())
			file_logger.Println("Invalid request payload:", err.Error())
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		if !common.ValidateAndRespond(w, movie) {
			return
		}
		for _, actor := range movie.Actors {
			if !common.ValidateAndRespond(w, actor) {
				return
			}
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

// ReadAllMovies returns a list of movies with pagination support.
// @Summary Get all movies with pagination
// @ID read-all-movies
// @Produce json
// @Security BearerAuth
// @Tags Movies
// @Param page query integer true "Page number"
// @Param page_size query integer true "Number of items per page"
// @Param sort_by query string false "Field to sort by (default rating)"
// @Param sort_order query string false "Sort order (ASC or DESC, default DESC)"
// @Success 200 {object} ReadAllMoviesResponse "List of movies"
// @Failure 400 {string} string "Invalid page_size or page"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/movies/all [get]
func ReadAllMovies(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.ReadAllMovieApiPingCounter.Inc()
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

// ReadMovie returns details of a specific movie by ID.
// @Summary Get movie by ID
// @ID read-movie-by-id
// @Produce json
// @Tags Movies
// @Security BearerAuth
// @Param id query string true "Movie ID"
// @Success 200 {object} ReadMovieResponse "Movie details"
// @Failure 400 {string} string "Movie ID is required"
// @Failure 404 {string} string "Movie not found"
// @Failure 500 {string} string "Failed to fetch movie"
// @Router /api/v1/movies/get [get]
func ReadMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.ReadOneMovieApiPingCounter.Inc()
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

// UpdateMovie updates an existing movie.
// @Summary Update an existing movie
// @ID update-movie
// @Accept json
// @Tags Movies
// @Produce json
// @Security BearerAuth
// @Param movie body models.UpdateMovieModel true "Movie object to be updated"
// @Success 200 {object} UpdateMovieResponse "Movie updated successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Failed to update movie"
// @Router /api/v1/movies/update [put]
func UpdateMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.UpdateMovieApiPingCounter.Inc()
		w.Header().Set("Content-Type", "application/json")
		var movie models.Movie
		var response models.Movie
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		if !common.ValidateAndRespond(w, movie) {
			return
		}
		for _, actor := range movie.Actors {
			if !common.ValidateAndRespond(w, actor) {
				return
			}
		}
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", movie.ID).Save(&movie).Error; err != nil {
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

// DeleteMovie deletes a movie and its associations from the database.
// @Summary Delete a movie
// @ID delete-movie
// @Tags Movies
// @Security BearerAuth
// @Produce json
// @Param id query string true "Movie ID"
// @Success 200 {object} DeleteMovieResponse "Movie deleted successfully"
// @Failure 400 {string} string "Movie ID is required"
// @Failure 500 {string} string "Failed to delete movie or its associations"
// @Router /api/v1/movies/delete [delete]
func DeleteMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.DeleteMovieApiPingCounter.Inc()
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

// PatchMovie updates an existing movie with the provided patch data.
// @Summary Update an existing movie partially
// @ID patch-movie
// @Accept json
// @Tags Movies
// @Produce json
// @Security BearerAuth
// @Param id query string true "Movie ID"
// @Param patchData body models.CreateMovieModel true "Patch data for updating the movie"
// @Success 200 {object} PatchMovieResponse "Movie updated successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Failed to update movie or its associations"
// @Router /api/v1/movies/patch [patch]
func PatchMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.PatchMovieApiPingCounter.Inc()
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
		actors, ok := patchData["actors"].([]interface{})
		if !ok {
			common.ErrorResponse(w, http.StatusBadRequest, "Actors data is missing or invalid")
			return
		}
		delete(patchData, "actors")
		tx_err := db.Transaction(func(tx *gorm.DB) error {
			for _, actor := range actors {
				actorMap, ok := actor.(map[string]interface{})
				if !ok {
					return errors.ErrUnsupported
				}
				if err := tx.Model(&models.Actor{}).Where("id=?", actorMap["id"]).Updates(actorMap).Error; err != nil {
					log.Error("Error updating actor:", err.Error())
					file_logger.Println("Error updating actor:", err.Error())
					return err
				}
				query := "SELECT COUNT(*) FROM actor_movies WHERE actor_id = ? AND movie_id = ?"
				var count int64
				if err := tx.Raw(query, actorMap["id"], movieID).Row().Scan(&count); err != nil {
					log.Error("Error counting actor-movie relationship:", err.Error())
					file_logger.Println("Error counting actor-movie relationship:", err.Error())
					return err
				}

				if count == 0 {
					query := "INSERT INTO actor_movies (actor_id, movie_id) VALUES (?, ?)"
					if err := tx.Exec(query, actorMap["id"], movieID).Error; err != nil {
						log.Error("Error inserting actor-movie relationship:", err.Error())
						file_logger.Println("Error inserting actor-movie relationship:", err.Error())
						return err
					}
				}
			}
			if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Model(&models.Movie{}).Where("id = ?", movieID).Updates(patchData).Error; err != nil {
				log.Error("Error updating actor:", err.Error())
				file_logger.Println("Error updating actor:", err.Error())
				return err
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
