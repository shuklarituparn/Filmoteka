package models

import (
	"time"
)

type Actor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FirstName string    `gorm:"size:100;not null" json:"first_name" validate:"required,max=100"`
	LastName  string    `gorm:"size:100;not null" json:"last_name" validate:"required,max=100"`
	Gender    string    `json:"gender" validate:"required"`
	BirthDate string    `gorm:"not null" json:"birth_date" validate:"required"`
	Movies    []Movie   `gorm:"many2many:actor_movies;" json:"movies,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
