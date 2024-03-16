package routes

import (
	"net/http"

	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/config"
	"github.com/shuklarituparn/Filmoteka/pkg/middleware"
)

func MovieRouter(mux *http.ServeMux) {
	const prefix = "/api/v1/movies"

	get_movie := prefix + "/get"
	mux.Handle(get_movie, middleware.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.ReadMovie(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	all_movies := prefix + "/all"
	mux.Handle(all_movies, middleware.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.ReadAllMovies(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	createRoute := prefix + "/create"
	mux.Handle(createRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.CreateMovie(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	updateRoute := prefix + "/update"
	mux.Handle(updateRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			controllers.UpdateMovie(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	patchRoute := prefix + "/patch"
	mux.Handle(patchRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			controllers.PatchMovie(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	deleteRoute := prefix + "/delete"
	mux.Handle(deleteRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			controllers.DeleteMovie(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))
}
