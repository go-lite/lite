package lite

import (
	"fmt"
	"net/http"
)

type Error interface {
	error
	StatusCode() int
}

type Violation struct {
	PropertyPath string         `json:"propertyPath,omitempty"`
	Message      string         `json:"message,omitempty"`
	Code         string         `json:"code,omitempty"`
	More         map[string]any `json:"more,omitempty"`
}

type HTTPError struct {
	Context     string      `json:"@context,omitempty"`
	Type        string      `json:"@type,omitempty"`
	Status      int         `json:"status,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Violations  []Violation `json:"violations,omitempty"`
}

func (e HTTPError) Error() string {
	title := e.Title

	if title == "" {
		title = http.StatusText(e.Status)

		if title == "" {
			title = "HTTP error response"
		}
	}

	return fmt.Sprintf("%s [%d]: %s", title, e.Status, e.Description)
}

func (e HTTPError) StatusCode() int {
	if e.Status == 0 {
		return StatusInternalServerError
	}

	return e.Status
}

func NewErrorResponse(status int, context, title, errorType string, description string, violations []Violation) HTTPError {
	return HTTPError{
		Context:     context,
		Type:        errorType,
		Status:      status,
		Title:       title,
		Description: description,
		Violations:  violations,
	}
}

func NewError(status int, descriptions ...string) HTTPError {
	if len(descriptions) == 0 {
		return NewErrorResponse(
			status,
			"/api/contexts/Error",
			"An error occurred",
			"Error",
			"An unexpected error occurred",
			nil,
		)
	}

	return NewErrorResponse(
		status,
		"/api/contexts/Error",
		"An error occurred",
		"Error",
		descriptions[0],
		nil,
	)
}

func NewNotFoundError(descriptions ...string) HTTPError {
	if len(descriptions) == 0 {
		return NewErrorResponse(
			StatusNotFound,
			"/api/contexts/NotFound",
			"Resource not found",
			"NotFound",
			"Resource not found",
			nil,
		)
	}

	return NewErrorResponse(
		StatusNotFound,
		"/api/contexts/NotFound",
		"Resource not found",
		"NotFound",
		descriptions[0],
		nil,
	)
}

func NewUnauthorizedError(descriptions ...string) HTTPError {
	if len(descriptions) == 0 {
		return NewErrorResponse(
			StatusUnauthorized,
			"/api/contexts/AuthenticationFailure",
			"Authentication failed",
			"AuthenticationFailure",
			"Authentication failed",
			nil,
		)
	}

	description := descriptions[0]

	return NewErrorResponse(
		StatusUnauthorized,
		"/api/contexts/AuthenticationFailure",
		"Authentication failed",
		"AuthenticationFailure",
		description,
		nil,
	)
}

func NewForbiddenError(descriptions ...string) HTTPError {
	if len(descriptions) == 0 {
		return NewErrorResponse(
			StatusForbidden,
			"/api/contexts/AccessDenied",
			"Access denied",
			"AccessDenied",
			"Access denied",
			nil,
		)
	}

	return NewErrorResponse(
		StatusForbidden,
		"/api/contexts/AccessDenied",
		"Access denied",
		"AccessDenied",
		descriptions[0],
		nil,
	)
}

func NewInternalServerError(descriptions ...string) HTTPError {
	if len(descriptions) == 0 {
		return NewErrorResponse(
			StatusInternalServerError,
			"/api/contexts/InternalServerError",
			"Internal server error",
			"InternalServerError",
			"Internal server error",
			nil,
		)
	}

	description := descriptions[0]

	return NewErrorResponse(
		StatusInternalServerError,
		"/api/contexts/InternalServerError",
		"Internal server error",
		"InternalServerError",
		description,
		nil,
	)
}

func NewServiceUnavailableError(descriptions ...string) HTTPError {
	if len(descriptions) == 0 {
		return NewErrorResponse(
			StatusServiceUnavailable,
			"/api/contexts/ServiceUnavailable",
			"Service unavailable",
			"ServiceUnavailable",
			"Service unavailable",
			nil,
		)
	}

	description := descriptions[0]

	return NewErrorResponse(
		StatusServiceUnavailable,
		"/api/contexts/ServiceUnavailable",
		"Service unavailable",
		"ServiceUnavailable",
		description,
		nil,
	)
}

func NewConflictError(descriptions ...string) HTTPError {
	if len(descriptions) == 0 {
		return NewErrorResponse(
			StatusConflict,
			"/api/contexts/Conflict",
			"Conflict",
			"Conflict",
			"Conflict",
			nil,
		)
	}

	description := descriptions[0]

	return NewErrorResponse(
		StatusConflict,
		"/api/contexts/Conflict",
		"Conflict",
		"Conflict",
		description,
		nil,
	)
}

func NewBadRequestError(descriptions ...string) HTTPError {
	if len(descriptions) == 0 {
		return NewErrorResponse(
			StatusBadRequest,
			"/api/contexts/BadRequest",
			"Bad request",
			"BadRequest",
			"Bad request",
			nil,
		)
	}

	description := descriptions[0]

	return NewErrorResponse(
		StatusBadRequest,
		"/api/contexts/BadRequest",
		"Bad request",
		"BadRequest",
		description,
		nil,
	)
}

func (e HTTPError) Descriptions() string {
	switch e.Status {
	case StatusBadRequest:
		return "Bad Request"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusForbidden:
		return "Forbidden"
	case StatusNotFound:
		return "Not Found"
	case StatusConflict:
		return "Conflict"
	case StatusInternalServerError:
		return "Internal Server Error"
	case StatusServiceUnavailable:
		return "Service Unavailable"
	default:
		return "Unknown Error"
	}
}

var DefaultErrorResponses = map[int]HTTPError{
	StatusBadRequest:          NewBadRequestError("Bad Request"),
	StatusInternalServerError: NewInternalServerError("Internal Server Error"),
}

var DefaultErrorContentTypeResponses = []string{
	"application/json",
	"application/xml",
	"multipart/form-data",
}

type BadRequestError HTTPError

func (e BadRequestError) Error() string {
	description := "Bad Request"

	if e.Description != "" {
		description = e.Description
	}

	return description
}

func (e BadRequestError) StatusCode() int {
	return StatusBadRequest
}

type UnauthorizedError HTTPError

func (e UnauthorizedError) Error() string {
	description := "Unauthorized"

	if e.Description != "" {
		description = e.Description
	}

	return description
}

func (e UnauthorizedError) StatusCode() int {
	return StatusUnauthorized
}

type ForbiddenError HTTPError

func (e ForbiddenError) Error() string {
	description := "Forbidden"

	if e.Description != "" {
		description = e.Description
	}

	return description
}

func (e ForbiddenError) StatusCode() int {
	return StatusForbidden
}

type NotFoundError HTTPError

func (e NotFoundError) Error() string {
	description := "Not Found"

	if e.Description != "" {
		description = e.Description
	}

	return description
}

func (e NotFoundError) StatusCode() int {
	return StatusNotFound
}

type ConflictError HTTPError

func (e ConflictError) Error() string {
	description := "Conflict"

	if e.Description != "" {
		description = e.Description
	}

	return description
}

func (e ConflictError) StatusCode() int {
	return StatusConflict
}

type InternalServerError HTTPError

func (e InternalServerError) Error() string {
	description := "Internal Server Error"

	if e.Description != "" {
		description = e.Description
	}

	return description
}

func (e InternalServerError) StatusCode() int {
	return StatusInternalServerError
}

type ServiceUnavailableError HTTPError

func (e ServiceUnavailableError) Error() string {
	description := "Service Unavailable"

	if e.Description != "" {
		description = e.Description
	}

	return description
}

func (e ServiceUnavailableError) StatusCode() int {
	return StatusServiceUnavailable
}
