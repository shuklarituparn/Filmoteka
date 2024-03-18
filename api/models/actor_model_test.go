package models

import (
	"github.com/go-playground/validator"
	"testing"
	"time"
)

func TestActorValidation(t *testing.T) {
	validate := validator.New()

	err := validate.RegisterValidation("validDate", ValidDateTest)
	if err != nil {
		return
	}

	validActor := Actor{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Gender:    "Male",
		BirthDate: "1990-01-01",
		Movies:    make([]Movie, 0),
	}
	if err := validate.Struct(validActor); err != nil {
		t.Errorf("Validation failed for valid actor: %v", err)
	}

	invalidActor := Actor{}

	if err := validate.Struct(invalidActor); err == nil {
		t.Error("Expected validation error for invalid actor, but got none")
	} else {
		t.Logf("Expected validation error for invalid actor: %v", err)
	}
}

func TestValidDateTest(t *testing.T) {
	validate := validator.New()
	err := validate.RegisterValidation("validDate", ValidDate)
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
}

func ValidDateTest(fl validator.FieldLevel) bool {
	birthDate, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	return birthDate.Before(time.Now())
}
