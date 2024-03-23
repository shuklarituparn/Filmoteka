package test

import (
	"github.com/go-playground/validator"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"testing"
)

func TestMovieValidation(t *testing.T) {

	r := runner.NewRunner(t, "Validation")
	r.NewTest("Movie model validation", func(P provider.T) {
		P.Feature("Movie model")
		P.Title("Testing movies model validation")
		P.Description("This test tests the movie model")
		validate := validator.New()
		validMovie := models.Movie{
			ID:          1,
			Title:       "Some Random Movie",
			ReleaseYear: 2000,
			Genre:       "Sci-fi",
			Description: "Alien Battles the humans",
			Rating:      10,
			Actors:      make([]models.Actor, 0),
		}
		if err := validate.Struct(validMovie); err != nil {
			t.Errorf("Validation failed for valid actor: %v", err)
		}

		invalidMovie := models.Movie{}

		if err := validate.Struct(invalidMovie); err == nil {
			t.Error("Expected validation error for invalid actor, but got none")
		} else {
			t.Logf("Expected validation error for invalid actor: %v", err)
		}
	})
	r.RunTests()

}
