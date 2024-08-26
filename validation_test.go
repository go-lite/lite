package lite

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Email string `validate:"required,email"`
	UUID  string `validate:"required,uuid"`
	Phone string `validate:"required,e164"`
}

func TestValidate_Success(t *testing.T) {
	app := &App{validator: validator.New()}
	testStruct := TestStruct{
		Email: "test@example.com",
		UUID:  uuid.NewString(),
		Phone: "+1234567890",
	}

	err := app.validate(testStruct)
	assert.NoError(t, err)
}

func TestValidate_RequiredFieldError(t *testing.T) {
	app := &App{validator: validator.New()}
	testStruct := TestStruct{
		Email: "test@example.com",
	}

	err := app.validate(testStruct)
	assert.Error(t, err)
}

func TestValidate_ValidationErrors(t *testing.T) {
	app := &App{validator: validator.New()}
	invalidStruct := TestStruct{
		Email: "invalid-email",
		UUID:  "invalid-uuid",
		Phone: "invalid-phone",
	}

	err := app.validate(invalidStruct)
	assert.Error(t, err)

	httpErr, ok := err.(HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Status)
	assert.Equal(t, "A constraint violation occurred", httpErr.Title)
	assert.Equal(t, "/api/contexts/ConstraintViolationList", httpErr.Context)

	expectedViolations := 3
	if len(httpErr.Violations) != expectedViolations {
		t.Errorf("expected %d violations, got %d", expectedViolations, len(httpErr.Violations))
	}

	expectedDescriptions := []string{
		"Email should be a valid email",
		"UUID should be a valid UUID",
		"Phone should be a valid international phone number (e.g. +33 1 23 45 67 89)",
	}
	actualDescriptions := strings.Split(httpErr.Description, ", ")
	for _, desc := range expectedDescriptions {
		assert.Contains(t, actualDescriptions, desc)
	}
}

func TestValidate_InvalidValidationError(t *testing.T) {
	app := &App{validator: validator.New()}

	// To simulate an InvalidValidationError, we need to pass an invalid type to validator.Struct
	invalidType := func() {}
	err := app.validate(invalidType)
	assert.Error(t, err)

	var invalidValidationError *validator.InvalidValidationError
	assert.True(t, errors.As(err, &invalidValidationError))
}

func TestValidate_MapInput(t *testing.T) {
	app := &App{validator: validator.New()}
	input := map[string]any{
		"key": "value",
	}

	err := app.validate(input)
	assert.NoError(t, err)
}
