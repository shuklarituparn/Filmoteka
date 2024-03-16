package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	"github.com/shuklarituparn/Filmoteka/pkg/hashing"
	jwt "github.com/shuklarituparn/Filmoteka/pkg/jwt_token"
	"gorm.io/gorm"
)

// RegisterUser handles user registration.
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.RegisterUserModel true "User information"
// @Success 201 {object} CreateUserResponse "User Created Successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/register [post]
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file_logger.Println("Request received:", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Error("Error decoding JSON:", err.Error())
			file_logger.Println("Error decoding JSON:", err.Error())
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		user.Role = "USER"
		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			errorsMap := make(map[string]string)
			for _, e := range err.(validator.ValidationErrors) {
				errorsMap[e.Field()] = e.Tag()
			}
			errJSON, _ := json.Marshal(errorsMap)
			common.ErrorResponse(w, http.StatusBadRequest, string(errJSON))
			return
		}
		hashedPassword, err := hashing.HashPassword(user.Password)
		if err != nil {
			common.ErrorResponse(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}
		user.Password = hashedPassword
		if err := db.Model(&models.User{}).Create(&user).Error; err != nil {
			log.Error("Error creating user:", err.Error())
			file_logger.Println("Error creating user:", err.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to register user %v", err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      user.ID,
			"message": "User Created Successfully",
			"email":   user.Email,
			"role":    user.Role,
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			file_logger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

// LoginUser handles user login.
// @Summary Log in a user
// @Description Log in a user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginUserModel true "User credentials"
// @Success 200 {object} LoginResponse "Logged In Successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 401 {object} ErrorResponse "Invalid email or password"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/login [post]
func LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file_logger.Println("Request received:", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		var storedUser models.User
		if err := db.Model(&models.User{}).Where("email = ?", user.Email).First(&storedUser).Error; err != nil {
			log.Error("Error finding user:", err.Error())
			file_logger.Println("Error finding user:", err.Error())
			common.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		if ok := hashing.CheckPasswordHash(user.Password, storedUser.Password); !ok {
			common.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		token, _ := jwt.GetJWTToken(user.Email, storedUser.Role, 1)
		refresh, err := jwt.GetJWTToken(user.Email, storedUser.Role, 5)
		if err != nil {
			log.Error("Error generating token:", err.Error())
			file_logger.Println("Error generating token:", err.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		resErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"id":            storedUser.ID,
			"message":       "Logged In Successfully",
			"email":         storedUser.Email,
			"role":          storedUser.Role,
			"access_token":  token,
			"refresh_token": refresh,
		})
		if resErr != nil {
			log.Error("Error encoding JSON:", resErr.Error())
			file_logger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

// RefreshToken handles refreshing JWT tokens.
// @Summary Refresh JWT tokens
// @Description Refresh JWT access and refresh tokens
// @Tags Authentication
// @Accept json
// @Security BearerAuth
// @Produce json
// @Success 200 {object} RefreshTokenResponse "New access and refresh tokens"
// @Failure 400 {object} ErrorResponse "Invalid or expired token"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /api/v1/users/refresh [get]
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	file_logger.Println("Request received:", r.Method, r.URL.Path)
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	if tokenString == "" {
		common.ErrorResponse(w, http.StatusBadRequest, "Please supply Token")
	}
	claims, err := jwt.VerifyToken(tokenString)
	if err != nil {
		common.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}
	newToken, _ := jwt.GetJWTToken(claims.Email, claims.Role, 1)
	refresh, err := jwt.GetJWTToken(claims.Email, claims.Role, 5)

	if err != nil {
		log.Error("Error generating new token:", err.Error())
		file_logger.Println("Error generating new token:", err.Error())
		common.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate new token")
		return
	}
	w.WriteHeader(http.StatusOK)
	resErr := json.NewEncoder(w).Encode(map[string]string{"token": newToken, "refresh": refresh})
	if resErr != nil {
		log.Error("Error encoding JSON:", resErr.Error())
		file_logger.Println("Error encoding JSON:", resErr.Error())
		common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
