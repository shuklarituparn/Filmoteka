package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	"gorm.io/gorm"
)

// SearchMovie searches for movies based on the provided query string.
// @Summary Search for movies
// @ID search-movies
// @Produce json
// @Tags Search Movies
// @Security BearerAuth
// @Param q query string true "Search query"
// @Param sort_by query string false "Field to sort by (default rating)"
// @Param sort_order query string false "Sort order (ASC or DESC, default DESC)"
// @Success 200 {object} SearchMovieResponse "List of matching movies"
// @Failure 400 {string} string "Invalid search query"
// @Failure 500 {string} string "Error encoding response"
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
			file_logger.Println("Error encoding response:", err.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Error encoding response")
			return
		}
	}
}
