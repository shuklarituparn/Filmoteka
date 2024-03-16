package models

import (
	"time"

	"github.com/go-playground/validator"
)

type Actor struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	FirstName string  `gorm:"size:100;not null" json:"first_name" validate:"required,max=100"`
	LastName  string  `gorm:"size:100;not null" json:"last_name" validate:"required,max=100"`
	Gender    string  `json:"gender" validate:"required"`
	BirthDate string  `gorm:"not null" json:"birth_date" validate:"required,validDate"`
	Movies    []Movie `gorm:"many2many:actor_movies;" json:"movies,omitempty"`
}

type CreateActorModel struct {
	FirstName string  `json:"first_name" validate:"required,max=100"`
	LastName  string  `json:"last_name" validate:"required,max=100"`
	Gender    string  `json:"gender" validate:"required"`
	BirthDate string  `json:"birth_date" validate:"required,validDate"`
	Movies    []Movie `json:"movies"`
}

type UpdateActorModel struct {
	CreateActorModel
	ID uint `json:"id" validate:"required"`
}

func ValidDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}
