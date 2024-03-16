package models

type Movie struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Title       string  `gorm:"size:150;not null" json:"title" validate:"required,min=1,max=150"`
	ReleaseYear int     `gorm:"not null" json:"release_year" validate:"required"`
	Genre       string  `gorm:"not null" json:"genre" validate:"required"`
	Description string  `gorm:"size:1000" json:"description" validate:"min=1,max=1000"`
	Rating      float64 `gorm:"not null" json:"rating" validate:"required,min=0,max=10"`
	Actors      []Actor `gorm:"many2many:actor_movies;" json:"actors,omitempty"`
}

type CreateMovieModel struct {
	Title       string  `json:"title" validate:"required,min=1,max=150"`
	ReleaseYear int     `json:"release_year"`
	Genre       string  `json:"genre"`
	Description string  `json:"description" validate:"max=1000"`
	Rating      float64 `json:"rating" validate:"required,min=0,max=10"`
	Actors      []Actor `json:"actors,omitempty"`
}

type UpdateMovieModel struct {
	CreateMovieModel
	ID uint `json:"id" validate:"required"`
}
