package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/internal/logger"
	"github.com/shuklarituparn/Filmoteka/internal/prometheus"
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	"gorm.io/gorm"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
)

var fileLogger = logger.SetupLogger()

// CreateActor создает нового актера.
// @Summary Создать нового актера
// @ID create-actor
// @Accept json
// @Produce json
// @Tags Актеры
// @Security BearerAuth
// @Param actor body models.CreateActorModel true "Объект актера, который нужно создать"
// @Success 201 {object} CreateActorResponse "Актер добавлен"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/v1/actors/create [post]
func CreateActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.CreateActorApiPingCounter.Inc()
		w.Header().Set("Content-Type", "application/json")
		var actor models.Actor
		if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "invalid request payload")
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "internal Server Error")
			}
		}(r.Body)
		if !common.ValidateAndRespond(w, actor) {
			return
		}
		for _, movie := range actor.Movies {
			if !common.ValidateAndRespond(w, movie) {
				return
			}
		}
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&actor).Error; err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to create actor")
				return err
			}
			return nil
		})
		if txErr != nil {
			return
		}
		w.WriteHeader(http.StatusCreated)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{"id": actor.ID, "message": "Actor Added"})
		if resErr != nil {
			common.ErrorResponse(w, http.StatusInternalServerError, "internal Server Error")
		}
	}
}

// ReadAllActor возвращает список актеров с поддержкой пагинации.
// @Summary Получить всех актеров с пагинацией
// @ID read-all-actors
// @Produce json
// @Tags Актеры
// @Param page query integer true "Номер страницы"
// @Param page_size query integer true "Количество элементов на странице"
// @Param sort_by query string false "Поле для сортировки (по умолчанию birth_date)"
// @Param sort_order query string false "Порядок сортировки (ASC или DESC, по умолчанию DESC)"
// @Success 200 {object} ReadAllActorResponse "Список актеров"
// @Failure 400 {string} string "Неверный размер страницы или номер страницы"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
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
			common.ErrorResponse(w, http.StatusBadRequest, "invalid page_size")
			return
		}
		pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || pageNum < 1 {
			common.ErrorResponse(w, http.StatusBadRequest, "invalid page")
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
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&models.Actor{}).Select("COUNT(*)").Count(&totalActorsCount).Error; err != nil {
				log.Error("Error counting actors:", err.Error())
				fileLogger.Println("Error counting actors:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
				return err
			}
			query := tx.Model(&actors).Limit(pageSize).Offset(offset).Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
			if err := query.Preload("Movies").Preload("Movies.Actors").Find(&actors).Error; err != nil {
				log.Error("Error fetching actors:", err.Error())
				fileLogger.Println("Error fetching actors:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
				return err
			}
			return nil
		})
		if txErr != nil {
			return
		}
		totalPages := int(math.Ceil(float64(totalActorsCount) / float64(pageSize)))
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": actors, "total_pages": totalPages}); err != nil {
			log.Error("Error encoding JSON:", err.Error())
			fileLogger.Println("Error encoding JSON:", err.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "internal Server Error")
		}
	}
}

// ReadActor возвращает подробности о конкретном актере по его идентификатору.
// @Summary Получить актера по идентификатору
// @ID read-actor-by-id
// @Produce json
// @Tags Актеры
// @Security BearerAuth
// @Param id query string true "Идентификатор актера"
// @Success 200 {object} ReadActorResponse "Подробности актера"
// @Failure 400 {string} string "Необходим идентификатор актера"
// @Failure 404 {string} string "Актер не найден"
// @Failure 500 {string} string "Ошибка при получении актера"
// @Router /api/v1/actors/get [get]
func ReadActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.ReadOneActorApiPingCounter.Inc()
		w.Header().Set("Content-Type", "application/json")
		var actor models.Actor
		actorID := r.URL.Query().Get("id")
		if actorID == "" {
			common.ErrorResponse(w, http.StatusBadRequest, "actor ID is required")
			return
		}
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Preload("Movies").Preload("Movies.Actors").First(&actor, actorID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					common.ErrorResponse(w, http.StatusNotFound, "actor not found")
					return err
				}
				log.Error("Error fetching actor:", err.Error())
				fileLogger.Println("Error fetching actor:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch actor")
				return err
			}
			return nil
		})
		if txErr != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{"data": actor})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			fileLogger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "internal Server Error")
		}
	}
}

// UpdateActor обновляет существующего актера.
// @Summary Обновить существующего актера
// @ID update-actor
// @Accept json
// @Security BearerAuth
// @Produce json
// @Tags Актеры
// @Param actor body models.UpdateActorModel true "Объект актера для обновления"
// @Success 200 {object} UpdateActorResponse "Актер успешно обновлен"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Не удалось обновить актера"
// @Router /api/v1/actors/update [put]
func UpdateActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.UpdateActorApiPingCounter.Inc()
		w.Header().Set("Content-Type", "application/json")
		var actor models.Actor
		var response models.Actor
		if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "invalid request payload")
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "internal Server Error")
			}
		}(r.Body)
		if !common.ValidateAndRespond(w, actor) {
			return
		}
		for _, movie := range actor.Movies {
			if !common.ValidateAndRespond(w, movie) {
				return
			}
		}
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Where("id = ?", actor.ID).Save(&actor).Error; err != nil {
				log.Error("Error updating actor:", err.Error())
				fileLogger.Println("Error updating actor:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "failed to update actor")
				return err
			}
			if err := tx.Preload("Movies").Preload("Movies.Actors").Find(&response, actor.ID).Error; err != nil {
				log.Error("Error updating actor:", err.Error())
				fileLogger.Println("Error updating actor:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "failed to update actor")
				return err
			}
			return nil
		})
		if txErr != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{"actor": response, "message": "actor updated successfully"})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			fileLogger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "internal Server Error")
		}
	}
}

// DeleteActor удаляет актера и его связи из базы данных.
// @Summary Удалить актера
// @ID delete-actor
// @Produce json
// @Tags Актеры
// @Security BearerAuth
// @Param id query string true "Идентификатор актера"
// @Success 200 {object} DeleteActorResponse "Актер успешно удален"
// @Failure 400 {string} string "Требуется идентификатор актера"
// @Failure 500 {string} string "Не удалось удалить актера или его связи"
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
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("DELETE FROM actor_movies WHERE actor_id = ?", actorID).Error; err != nil {
				log.Error("Error deleting association:", err.Error())
				fileLogger.Println("Error deleting association:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete Association")
				return err
			}
			if err := tx.Delete(&models.Actor{}, actorID).Error; err != nil {
				log.Error("Error deleting actor:", err.Error())
				fileLogger.Println("Error deleting actor:", err.Error())
				common.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete actor")
				return err
			}
			return nil
		})
		if txErr != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Actor deleted successfully",
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			fileLogger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "internal Server Error")
		}
	}
}

// PatchActor обновляет существующего актера предоставленными изменениями.
// @Summary Частично обновить существующего актера
// @ID patch-actor
// @Accept json
// @Security BearerAuth
// @Produce json
// @Tags Актеры
// @Param id query string true "Идентификатор актера"
// @Param patchData body models.CreateActorModel true "Данные для частичного обновления актера"
// @Success 200 {object} PatchActorResponse "Актер успешно обновлен"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Не удалось обновить актера или его связи"
// @Router /api/v1/actors/patch [patch]
func PatchActor(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prometheus.PatchActorApiPingCounter.Inc()
		w.Header().Set("Content-Type", "application/json")
		var patchData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&patchData); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "invalid request payload")
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "internal Server Error")
			}
		}(r.Body)
		actorID := r.URL.Query().Get("id")
		if actorID == "" {
			common.ErrorResponse(w, http.StatusBadRequest, "Actor ID is required")
			return
		}
		txErr := db.Transaction(func(tx *gorm.DB) error {
			if moviesData, ok := patchData["movies"]; ok {
				movies, ok := moviesData.([]interface{})
				if !ok {
					common.ErrorResponse(w, http.StatusBadRequest, "Movies data is invalid")
				}
				delete(patchData, "movies")
				for _, movie := range movies {
					movieMap, ok := movie.(map[string]interface{})
					if !ok {
						return errors.ErrUnsupported
					}
					if err := tx.Model(&models.Movie{}).Where("id=?", movieMap["id"]).Updates(movieMap).Error; err != nil {
						log.Error("Error updating movie:", err.Error())
						fileLogger.Println("Error updating movie:", err.Error())
						return err
					}
					query := "SELECT COUNT(*) FROM actor_movies WHERE actor_id = ? AND movie_id = ?"
					var count int64
					if err := tx.Raw(query, actorID, movieMap["id"]).Row().Scan(&count); err != nil {
						log.Error("Error counting actor-movie relationship:", err.Error())
						fileLogger.Println("Error counting actor-movie relationship:", err.Error())
						return err
					}
					if count == 0 {
						query := "INSERT INTO actor_movies (actor_id, movie_id) VALUES (?, ?)"
						if err := tx.Exec(query, actorID, movieMap["id"]).Error; err != nil {
							log.Error("Error inserting actor-movie relationship:", err.Error())
							fileLogger.Println("Error inserting actor-movie relationship:", err.Error())
							return err
						}
					}
				}
			}
			if err := tx.Model(&models.Actor{}).Where("id = ?", actorID).Updates(patchData).Error; err != nil {
				log.Error("Error updating actor:", err.Error())
				fileLogger.Println("Error updating actor:", err.Error())
				common.ErrorResponse(w, http.StatusBadRequest, "Error updating actor")
				return err
			}
			return nil
		})
		if txErr != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Actor updated successfully",
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			fileLogger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "internal Server Error")
		}
	}
}
