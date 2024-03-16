package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"password" validate:"required,min=8"`
	Role     string `gorm:"not null" json:"role" validate:"required,oneof=USER ADMIN"`
}

type RegisterUserModel struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUserModel struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
