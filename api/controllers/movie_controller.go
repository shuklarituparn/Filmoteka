package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/internal/prometheus"
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	"gorm.io/gorm"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// CreateMovie создает новый фильм.
// @Summary Создать новый фильм
// @ID create-movie
// @Security BearerAuth
// @Accept json
// @Produce json
// @Tags Фильмы
// @Param movie body models.CreateMovieModel true "Объект фильма для создания"
// @Success 201 {object} CreateMovieResponse "Фильм успешно создан"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/v1/movies/create [post]
func CreateMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.CreateMovieApiPingCounter.Inc()
		w.Header().Set("Content-Type", "application/json")
		var movie models.Movie
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			log.Error("Invalid request payload:", err.Error())
			fileLogger.Println("Invalid request payload:", err.Error())
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}(r.Body)
		if !common.ValidateAndRespond(w, movie) {
			return
		}
		for _, actor := range movie.Actors {
			if !common.ValidateAndRespond(w, actor) {
				return
			}
		}
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&movie).Error; err != nil {
				log.Error("Failed to create movie:", err.Error())
				fileLogger.Println("Failed to create movie:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("%v", err.Error()))
				return err
			}
			return nil
		})
		if txErr != nil {
			return
		}
		w.WriteHeader(http.StatusCreated)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Movie created successfully",
			"data":    movie,
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			fileLogger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

// ReadAllMovies возвращает список фильмов с поддержкой пагинации.
// @Summary Получить все фильмы с пагинацией
// @ID read-all-movies
// @Produce json
// @Security BearerAuth
// @Tags Фильмы
// @Param page query integer true "Номер страницы"
// @Param page_size query integer true "Количество элементов на странице"
// @Param sort_by query string false "Поле для сортировки (по умолчанию рейтинг)"
// @Param sort_order query string false "Порядок сортировки (ASC или DESC, по умолчанию DESC)"
// @Success 200 {object} ReadAllMoviesResponse "Список фильмов"
// @Failure 400 {string} string "Неверный размер страницы или номер страницы"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
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
		txErr := db.Transaction(func(tx *gorm.DB) error {
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
		if txErr != nil {
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

// ReadMovie возвращает подробности о конкретном фильме по его идентификатору.
// @Summary Получить фильм по идентификатору
// @ID read-movie-by-id
// @Produce json
// @Tags Фильмы
// @Security BearerAuth
// @Param id query string true "Идентификатор фильма"
// @Success 200 {object} ReadMovieResponse "Подробности фильма"
// @Failure 400 {string} string "Требуется идентификатор фильма"
// @Failure 404 {string} string "Фильм не найден"
// @Failure 500 {string} string "Ошибка при получении фильма"
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
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Preload("Actors").Preload("Actors.Movies").First(&movie, movieID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					common.ErrorResponse(w, http.StatusNotFound, "Movie not found")
					return err
				}
				log.Error("Error fetching movie:", err.Error())
				fileLogger.Println("Error fetching movie:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
				return err
			}
			return nil
		})
		if txErr != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"data": movie,
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			fileLogger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

// UpdateMovie обновляет существующий фильм.
// @Summary Обновить существующий фильм
// @ID update-movie
// @Accept json
// @Tags Фильмы
// @Produce json
// @Security BearerAuth
// @Param movie body models.UpdateMovieModel true "Объект фильма для обновления"
// @Success 200 {object} UpdateMovieResponse "Фильм успешно обновлен"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Не удалось обновить фильм"
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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}(r.Body)
		if !common.ValidateAndRespond(w, movie) {
			return
		}
		for _, actor := range movie.Actors {
			if !common.ValidateAndRespond(w, actor) {
				return
			}
		}
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", movie.ID).Save(&movie).Error; err != nil {
				log.Error("Failed to update movie:", err.Error())
				fileLogger.Println("Failed to update movie:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to update movie")
				return err
			}
			if err := tx.Preload("Actors").Preload("Actors.Movies").Find(&response, movie.ID).Error; err != nil {
				log.Error("Failed to update movie:", err.Error())
				fileLogger.Println("Failed to update movie:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to update movie")
				return err
			}
			return nil
		})
		if txErr != nil {
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

// DeleteMovie удаляет фильм и его связи из базы данных.
// @Summary Удалить фильм
// @ID delete-movie
// @Tags Фильмы
// @Security BearerAuth
// @Produce json
// @Param id query string true "Идентификатор фильма"
// @Success 200 {object} DeleteMovieResponse "Фильм успешно удален"
// @Failure 400 {string} string "Требуется идентификатор фильма"
// @Failure 500 {string} string "Не удалось удалить фильм или его связи"
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
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("DELETE FROM actor_movies WHERE movie_id = ?", movieID).Error; err != nil {
				log.Error("Failed to delete Association:", err.Error())
				fileLogger.Println("Failed to delete Association:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete Association")
				return err
			}
			if err := tx.Delete(&models.Movie{}, movieID).Error; err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete movie")
				return err
			}
			return nil
		})
		if txErr != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Movie deleted successfully",
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			fileLogger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

// PatchMovie обновляет существующий фильм предоставленными изменениями.
// @Summary Частично обновить существующий фильм
// @ID patch-movie
// @Accept json
// @Tags Фильмы
// @Produce json
// @Security BearerAuth
// @Param id query string true "Идентификатор фильма"
// @Param patchData body models.CreateMovieModel true "Данные для частичного обновления фильма"
// @Success 200 {object} PatchMovieResponse "Фильм успешно обновлен"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Не удалось обновить фильм или его связи"
// @Router /api/v1/movies/patch [patch]
func PatchMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.PatchMovieApiPingCounter.Inc()
		w.Header().Set("Content-Type", "application/json")
		var patchData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
			log.Error("Invalid request payload:", err.Error())
			fileLogger.Println("Invalid request payload:", err.Error())
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}(r.Body)
		movieID := r.URL.Query().Get("id")
		if movieID == "" {
			common.ErrorResponse(w, http.StatusBadRequest, "Movie ID is required")
			return
		}
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if actorsData, ok := patchData["actors"]; ok {
				actors, ok := actorsData.([]interface{})
				if !ok {
					common.ErrorResponse(w, http.StatusBadRequest, "Movies data is invalid")
				}
				delete(patchData, "actors")
				for _, actor := range actors {
					actorMap, ok := actor.(map[string]interface{})
					if !ok {
						return errors.ErrUnsupported
					}
					if err := tx.Model(&models.Actor{}).Where("id=?", actorMap["id"]).Updates(actorMap).Error; err != nil {
						log.Error("Error updating actor:", err.Error())
						fileLogger.Println("Error updating actor:", err.Error())
						return err
					}
					query := "SELECT COUNT(*) FROM actor_movies WHERE actor_id = ? AND movie_id = ?"
					var count int64
					if err := tx.Raw(query, actorMap["id"], movieID).Row().Scan(&count); err != nil {
						log.Error("Error counting actor-movie relationship:", err.Error())
						fileLogger.Println("Error counting actor-movie relationship:", err.Error())
						return err
					}

					if count == 0 {
						query := "INSERT INTO actor_movies (actor_id, movie_id) VALUES (?, ?)"
						if err := tx.Exec(query, actorMap["id"], movieID).Error; err != nil {
							log.Error("Error inserting actor-movie relationship:", err.Error())
							fileLogger.Println("Error inserting actor-movie relationship:", err.Error())
							return err
						}
					}
				}
			}
			if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Model(&models.Movie{}).Where("id = ?", movieID).Updates(patchData).Error; err != nil {
				log.Error("Error updating movie:", err.Error())
				fileLogger.Println("Error updating movie:", err.Error())
				common.ErrorResponse(w, http.StatusBadRequest, "Error updating movie")
				return err
			}
			return nil
		})
		if txErr != nil {
			return
		}
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Movie updated successfully",
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			fileLogger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}
