package models

import (
	"github.com/go-playground/validator"
	"testing"
)

func TestMovieValidation(t *testing.T) {
	validate := validator.New()
	validMovie := Movie{
		ID:          1,
		Title:       "Some Random Movie",
		ReleaseYear: 2000,
		Genre:       "Sci-fi",
		Description: "Alien Battles the humans",
		Rating:      10,
		Actors:      make([]Actor, 0),
	}
	if err := validate.Struct(validMovie); err != nil {
		t.Errorf("Validation failed for valid actor: %v", err)
	}

	invalidMovie := Movie{}

	if err := validate.Struct(invalidMovie); err == nil {
		t.Error("Expected validation error for invalid actor, but got none")
	} else {
		t.Logf("Expected validation error for invalid actor: %v", err)
	}
}
