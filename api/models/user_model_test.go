package models

import (
	"github.com/go-playground/validator"
	"testing"
)

func TestUserValidation(t *testing.T) {
	validate := validator.New()
	validUser := User{
		ID:       1,
		Email:    "some@example.com",
		Password: "password",
		Role:     "USER",
	}
	if err := validate.Struct(validUser); err != nil {
		t.Errorf("Validation failed for valid actor: %v", err)
	}

	invalidUser := User{}

	if err := validate.Struct(invalidUser); err == nil {
		t.Error("Expected validation error for invalid actor, but got none")
	} else {
		t.Logf("Expected validation error for invalid actor: %v", err)
	}
}
