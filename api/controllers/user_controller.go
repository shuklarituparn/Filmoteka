package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/go-playground/validator"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"github.com/shuklarituparn/Filmoteka/pkg/common"
	"github.com/shuklarituparn/Filmoteka/pkg/hashing"
	jwt "github.com/shuklarituparn/Filmoteka/pkg/jwt_token"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strings"
)

// RegisterUser обрабатывает регистрацию пользователя.
// @Summary Зарегистрировать нового пользователя
// @Description Зарегистрировать нового пользователя с помощью электронной почты и пароля
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body models.RegisterUserModel true "Информация о пользователе"
// @Success 201 {object} CreateUserResponse "Пользователь успешно создан"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/register [post]
func RegisterUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileLogger.Println("Request received:", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Error("Error decoding JSON:", err.Error())
			fileLogger.Println("Error decoding JSON:", err.Error())
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}(r.Body)
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
			fileLogger.Println("Error creating user:", err.Error())
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
			fileLogger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

// LoginUser обрабатывает вход пользователя.
// @Summary Вход пользователя
// @Description Вход пользователя с использованием электронной почты и пароля
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body models.LoginUserModel true "Учетные данные пользователя"
// @Success 200 {object} LoginResponse "Успешный вход в систему"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 401 {object} ErrorResponse "Неверная электронная почта или пароль"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/login [post]
func LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileLogger.Println("Request received:", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			common.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}(r.Body)
		var storedUser models.User
		if err := db.Model(&models.User{}).Where("email = ?", user.Email).First(&storedUser).Error; err != nil {
			log.Error("Error finding user:", err.Error())
			fileLogger.Println("Error finding user:", err.Error())
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
			fileLogger.Println("Error generating token:", err.Error())
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
			fileLogger.Println("Error encoding JSON:", resErr.Error())
			common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

// RefreshToken обрабатывает обновление токенов JWT.
// @Summary Обновление токенов JWT
// @Description Обновление доступа и обновления токенов JWT
// @Tags Аутентификация
// @Accept json
// @Security BearerAuth
// @Produce json
// @Success 200 {object} RefreshTokenResponse "Новые токены доступа и обновления"
// @Failure 400 {object} ErrorResponse "Неверный или истекший токен"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/users/refresh [get]
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fileLogger.Println("Request received:", r.Method, r.URL.Path)
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
		fileLogger.Println("Error generating new token:", err.Error())
		common.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate new token")
		return
	}
	w.WriteHeader(http.StatusOK)
	resErr := json.NewEncoder(w).Encode(map[string]string{"token": newToken, "refresh": refresh})
	if resErr != nil {
		log.Error("Error encoding JSON:", resErr.Error())
		fileLogger.Println("Error encoding JSON:", resErr.Error())
		common.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
