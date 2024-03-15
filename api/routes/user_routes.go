package routes

import (
	"net/http"

	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/config"
	"github.com/shuklarituparn/Filmoteka/pkg/middleware"
)

func UserRouter(mux *http.ServeMux) {
	const prefix = "/api/v1/users"
	registerRoute := prefix + "/register"
	loginRoute := prefix + "/login"
	refreshTokenRoute := prefix + "/refresh"

	mux.Handle(registerRoute, middleware.IsAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.RegisterUser(config.GetInstance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle(loginRoute, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.LoginUser(config.GetInstance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.Handle(refreshTokenRoute, middleware.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.RefreshToken(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))
}
