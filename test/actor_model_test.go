package test

import (
	"github.com/go-playground/validator"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/shuklarituparn/Filmoteka/api/models"
	"testing"
	"time"
)

func TestActorValidation(t *testing.T) {

	r := runner.NewRunner(t, "Validation")
	r.NewTest("Actor model validation", func(P provider.T) {
		P.Feature("Actor model")
		P.Title("Testing actor model validation")
		P.Description("This test the actor model")
		validate := validator.New()

		err := validate.RegisterValidation("validDate", ValidDateTest)
		if err != nil {
			return
		}

		validActor := models.Actor{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Gender:    "Male",
			BirthDate: "1990-01-01",
			Movies:    make([]models.Movie, 0),
		}
		if err := validate.Struct(validActor); err != nil {
			t.Errorf("Validation failed for valid actor: %v", err)
		}

		invalidActor := models.Actor{}

		if err := validate.Struct(invalidActor); err == nil {
			t.Error("Expected validation error for invalid actor, but got none")
		} else {
			t.Logf("Expected validation error for invalid actor: %v", err)
		}

	})
	r.RunTests()

}

func TestValidDateTest(t *testing.T) {

	r := runner.NewRunner(t, "Validation")
	r.NewTest("Date validation", func(P provider.T) {
		P.Feature("Date model")
		P.Title("Checking if the date validator work")
		P.Description("This test checks if the validation works")
		validate := validator.New()
		err := validate.RegisterValidation("validDate", models.ValidDate)
		if err != nil {
			return
		}
		type TestStruct struct {
			Date string `validate:"required,validDate"`
		}

		cases := []struct {
			name     string
			dateStr  string
			expected bool
		}{
			{"Valid date", "2024-03-17", true},
			{"Invalid date format", "17-03-2024", false},
			{"Invalid date value", "2024-02-31", false},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				testStruct := TestStruct{Date: tc.dateStr}
				err := validate.Struct(testStruct)

				result := err == nil
				if result != tc.expected {
					t.Errorf("Expected %v for date %s; Got %v", tc.expected, tc.dateStr, result)
				}
			})
		}

	})
	r.RunTests()

}

func ValidDateTest(fl validator.FieldLevel) bool {
	birthDate, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	return birthDate.Before(time.Now())
}
