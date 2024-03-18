package routes

import (
	"github.com/shuklarituparn/Filmoteka/config"
	"net/http"

	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/pkg/middleware"
)

func ActorRouter(mux *http.ServeMux) {
	const prefix = "/api/v1/actors"

	getActor := prefix + "/get"
	mux.Handle(getActor, middleware.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.ReadActor(config.GetInstance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	allActors := prefix + "/all"
	mux.Handle(allActors, middleware.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.ReadAllActor(config.GetInstance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	updateRoute := prefix + "/update"
	mux.Handle(updateRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			controllers.UpdateActor(config.GetInstance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	patchRoute := prefix + "/patch"
	mux.Handle(patchRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			controllers.PatchActor(config.GetInstance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	deleteRoute := prefix + "/delete"
	mux.Handle(deleteRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			controllers.DeleteActor(config.GetInstance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	createRoute := prefix + "/create"
	mux.Handle(createRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.CreateActor(config.GetInstance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))
}
