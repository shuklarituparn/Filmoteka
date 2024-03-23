package test

import (
	"github.com/go-playground/validator"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"testing"
)

func TestUserValidation(t *testing.T) {

	r := runner.NewRunner(t, "Validation")
	r.NewTest("User model validation", func(P provider.T) {
		P.Feature("User model")
		P.Title("Testing user model validation")
		P.Description("This test tests the user model")
		validate := validator.New()
		validUser := models.User{
			ID:       1,
			Email:    "some@example.com",
			Password: "password",
			Role:     "USER",
		}
		if err := validate.Struct(validUser); err != nil {
			t.Errorf("Validation failed for valid actor: %v", err)
		}

		invalidUser := models.User{}

		if err := validate.Struct(invalidUser); err == nil {
			t.Error("Expected validation error for invalid actor, but got none")
		} else {
			t.Logf("Expected validation error for invalid actor: %v", err)
		}
	})
	r.RunTests()

}
