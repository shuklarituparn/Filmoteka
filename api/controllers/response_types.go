package controllers

import (
	"time"

	"github.com/shuklarituparn/Filmoteka/api/models"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateActorResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type ReadAllActorResponse struct {
	Data       []models.Actor `json:"data"`
	TotalPages int            `json:"total_pages"`
}

type ReadActorResponse struct {
	Data models.Actor `json:"data"`
}

type UpdateActorResponse struct {
	Actor   models.Actor `json:"actor"`
	Message string       `json:"message"`
}

type DeleteActorResponse struct {
	Message string `json:"message"`
}

type PatchActorResponse struct {
	Message string `json:"message"`
}

type HealthCheckResponse struct {
	Author      string    `json:"author"`
	CurrentTime time.Time `json:"current_time"`
	Status      string    `json:"status"`
}

type CreateMovieResponse struct {
	Message string       `json:"message"`
	Data    models.Movie `json:"data"`
}

type ReadAllMoviesResponse struct {
	Data       []models.Movie `json:"data"`
	TotalPages int            `json:"total_pages"`
}

type ReadMovieResponse struct {
	Data models.Movie `json:"data"`
}

type UpdateMovieResponse struct {
	Message string       `json:"message"`
	Data    models.Movie `json:"data"`
}

type DeleteMovieResponse struct {
	Message string `json:"message"`
}

type PatchMovieResponse struct {
	Message string `json:"message"`
}

type SearchMovieResponse struct {
	Data []models.Movie `json:"data"`
}

type RefreshTokenResponse struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	Email        string `json:"email"`
	ID           int    `json:"id"`
	Message      string `json:"message"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

type CreateUserResponse struct {
	Email   string `json:"email"`
	ID      int    `json:"id"`
	Message string `json:"message"`
	Role    string `json:"role"`
}
