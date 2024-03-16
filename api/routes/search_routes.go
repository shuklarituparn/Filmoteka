package routes

import (
	"net/http"

	"github.com/shuklarituparn/Filmoteka/api/controllers"
	"github.com/shuklarituparn/Filmoteka/config"
	"github.com/shuklarituparn/Filmoteka/pkg/middleware"
)

func SearchRouter(mux *http.ServeMux) {
	const prefix = "/api/v1/search"

	mux.Handle(prefix, middleware.IsAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			controllers.SearchMovie(config.Get_Instance())(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})))
}
