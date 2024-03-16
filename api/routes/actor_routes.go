package routes

import (
	"net/http"

	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/config"
	"github.com/shuklarituparn/Filmoteka/pkg/middleware"
)

func ActorRouter(mux *http.ServeMux) {
	const prefix = "/api/v1/actors"

	get_actor := prefix + "/get"
	mux.Handle(get_actor, middleware.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.ReadActor(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	all_actors := prefix + "/all"
	mux.Handle(all_actors, middleware.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.ReadAllActor(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	updateRoute := prefix + "/update"
	mux.Handle(updateRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			controllers.UpdateActor(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	patchRoute := prefix + "/patch"
	mux.Handle(patchRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			controllers.PatchActor(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	deleteRoute := prefix + "/delete"
	mux.Handle(deleteRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			controllers.DeleteActor(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	createRoute := prefix + "/create"
	mux.Handle(createRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.CreateActor(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))
}
