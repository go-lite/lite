package lite

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"
)

// describeError translates a validator error into a human readable string.
func describeError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "uuid":
		return err.Field() + " should be a valid UUID"
	case "email":
		return err.Field() + " should be a valid email"
	case "e164":
		return err.Field() + " should be a valid international phone number (e.g. +33 1 23 45 67 89)"
	default:
		resp := fmt.Sprintf("%s should be %s", err.Field(), err.Tag())
		if err.Param() != "" {
			resp += "=" + err.Param()
		}

		return resp
	}
}

func (s *App) validate(a any) error {
	_, ok := a.(map[string]any)
	if ok {
		return nil
	}

	err := s.validator.Struct(a)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return fmt.Errorf("validation error: %w", err)
		}

		validationError := HTTPError{
			Context: "/api/contexts/ConstraintViolationList",
			Type:    "ConstraintViolation",
			Status:  http.StatusBadRequest,
			Title:   "A constraint violation occurred",
		}

		var errorsDescription []string
		var validationErrs validator.ValidationErrors

		if errors.As(err, &validationErrs) {
			for _, err := range validationErrs {
				errorsDescription = append(errorsDescription, describeError(err))
				validationError.Violations = append(validationError.Violations, Violation{
					PropertyPath: err.Field(),
					Message:      err.Error(),
					Code:         uuid.NewString(),
				})
			}
		}

		validationError.Description = strings.Join(errorsDescription, ", ")

		return validationError
	}

	return nil
}
