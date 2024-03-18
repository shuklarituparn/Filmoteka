package middleware

import (
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	jwt "github.com/shuklarituparn/Filmoteka/pkg/jwt_token"
)

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("Request received:", r.Method, r.URL.Path)
		fileLogger.Println("Request received:", r.Method, r.URL.Path)
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			common.ErrorResponse(w, http.StatusUnauthorized, "Authorization header missing")
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid authorization header format")
			return
		}
		token := tokenParts[1]
		_, err := jwt.VerifyToken(token)
		if err != nil {
			common.ErrorResponse(w, http.StatusForbidden, "Token Not Valid")
			return
		}
		next.ServeHTTP(w, r)
	})
}
