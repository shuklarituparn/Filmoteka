package models

import (
	"time"
)

type Movie struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:150;not null" json:"title" validate:"required,min=1,max=150"`
	ReleaseYear int       `gorm:"not null" json:"release_year"`
	Genre       string    `gorm:"not null" json:"genre"`
	Description string    `gorm:"size:1000" json:"description" validate:"max=1000"`
	Rating      float64   `gorm:"not null" json:"rating" validate:"required,min=0,max=10"`
	Actors      []Actor   `gorm:"many2many:actor_movies;" json:"actors,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
