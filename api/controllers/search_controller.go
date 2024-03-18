package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// SearchMovie выполняет поиск фильмов на основе предоставленной строки запроса.
// @Summary Поиск фильмов
// @ID search-movies
// @Produce json
// @Tags Поиск Фильмов
// @Security BearerAuth
// @Param q query string true "Поисковый запрос"
// @Param sort_by query string false "Поле для сортировки (по умолчанию рейтинг)"
// @Param sort_order query string false "Порядок сортировки (ASC или DESC, по умолчанию DESC)"
// @Success 200 {object} SearchMovieResponse "Список совпадающих фильмов"
// @Failure 400 {string} string "Неверный поисковый запрос"
// @Failure 500 {string} string "Ошибка при кодировании ответа"
// @Router /api/v1/search [get]
func SearchMovie(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		query := r.URL.Query().Get("q")
		sortBy := r.URL.Query().Get("sort_by")
		sortOrder := strings.ToUpper(r.URL.Query().Get("sort_order"))
		if sortBy == "" {
			sortBy = "rating"
		}
		if sortOrder != "ASC" && sortOrder != "DESC" {
			sortOrder = "DESC"
		}
		if query == "" {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid search query")
			return
		}
		searchQuery := "%" + strings.TrimSpace(query) + "%"

		var movies []models.Movie
		db.Model(&models.Movie{}).
			Joins("JOIN actor_movies ON movies.id = actor_movies.movie_id").
			Joins("JOIN actors ON actors.id = actor_movies.actor_id").
			Where("LOWER(title) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?", "%"+strings.ToLower(searchQuery)+"%", "%"+strings.ToLower(searchQuery)+"%", "%"+strings.ToLower(searchQuery)+"%").
			Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
			Preload("Actors").Preload("Actors.Movies").
			Distinct().
			Find(&movies)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": movies}); err != nil {
			log.Error("Error encoding response:", err.Error())
			fileLogger.Println("Error encoding response:", err.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Error encoding response")
			return
		}
	}
}
