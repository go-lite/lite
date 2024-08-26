package lite

import (
	"net/http"
	"testing"
)

func TestErrorResponse_Error(t *testing.T) {
	tests := []struct {
		name     string
		response HTTPError
		expected string
	}{
		{
			name: "Error with title",
			response: HTTPError{
				Status:      http.StatusNotFound,
				Title:       "Not Found",
				Description: "Resource not found",
			},
			expected: "Not Found [404]: Resource not found",
		},
		{
			name: "Error without title but with status",
			response: HTTPError{
				Status:      http.StatusInternalServerError,
				Description: "Server error",
			},
			expected: "Internal Server Error [500]: Server error",
		},
		{
			name: "Error without title and without known status",
			response: HTTPError{
				Status:      999,
				Description: "Unknown error",
			},
			expected: "HTTP error response [999]: Unknown error",
		},
		{
			name: "Error with empty status",
			response: HTTPError{
				Status:      0,
				Description: "Empty status error",
			},
			expected: "HTTP error response [0]: Empty status error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.Error(); got != tt.expected {
				t.Errorf("ErrorResponse.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorResponse_StatusCode(t *testing.T) {
	tests := []struct {
		name     string
		response HTTPError
		expected int
	}{
		{
			name: "Status code provided",
			response: HTTPError{
				Status: http.StatusBadRequest,
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Status code not provided, default to internal server error",
			response: HTTPError{
				Status: 0,
			},
			expected: StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.StatusCode(); got != tt.expected {
				t.Errorf("ErrorResponse.StatusCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewErrorResponse(t *testing.T) {
	violations := []Violation{
		{
			PropertyPath: "field1",
			Message:      "error message",
			Code:         "123",
		},
	}

	response := NewErrorResponse(http.StatusBadRequest, "/context", "Bad Request", "Error", "This is a bad request", violations)

	if response.Status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, response.Status)
	}

	if response.Title != "Bad Request" {
		t.Errorf("Expected title 'Bad Request', got '%s'", response.Title)
	}

	if len(response.Violations) != 1 || response.Violations[0].PropertyPath != "field1" {
		t.Errorf("Expected one violation with PropertyPath 'field1', got %v", response.Violations)
	}
}

func TestNewError(t *testing.T) {
	response := NewError(http.StatusInternalServerError)

	if response.Status != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, response.Status)
	}

	if response.Description != "An unexpected error occurred" {
		t.Errorf("Expected description 'An unexpected error occurred', got '%s'", response.Description)
	}

	customResponse := NewError(http.StatusInternalServerError, "Custom error")
	if customResponse.Description != "Custom error" {
		t.Errorf("Expected custom description 'Custom error', got '%s'", customResponse.Description)
	}
}

func TestNewNotFoundError(t *testing.T) {
	response := NewNotFoundError()

	if response.Status != StatusNotFound {
		t.Errorf("Expected status %d, got %d", StatusNotFound, response.Status)
	}

	if response.Description != "Resource not found" {
		t.Errorf("Expected description 'Resource not found', got '%s'", response.Description)
	}

	customResponse := NewNotFoundError("Custom not found error")
	if customResponse.Description != "Custom not found error" {
		t.Errorf("Expected custom description 'Custom not found error', got '%s'", customResponse.Description)
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	response := NewUnauthorizedError()

	if response.Status != StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", StatusUnauthorized, response.Status)
	}

	if response.Description != "Authentication failed" {
		t.Errorf("Expected description 'Authentication failed', got '%s'", response.Description)
	}

	customResponse := NewUnauthorizedError("Custom unauthorized error")
	if customResponse.Description != "Custom unauthorized error" {
		t.Errorf("Expected custom description 'Custom unauthorized error', got '%s'", customResponse.Description)
	}
}

func TestNewForbiddenError(t *testing.T) {
	response := NewForbiddenError()

	if response.Status != StatusForbidden {
		t.Errorf("Expected status %d, got %d", StatusForbidden, response.Status)
	}

	if response.Description != "Access denied" {
		t.Errorf("Expected description 'Access denied', got '%s'", response.Description)
	}

	customResponse := NewForbiddenError("Custom forbidden error")
	if customResponse.Description != "Custom forbidden error" {
		t.Errorf("Expected custom description 'Custom forbidden error', got '%s'", customResponse.Description)
	}
}

func TestNewInternalServerError(t *testing.T) {
	response := NewInternalServerError()

	if response.Status != StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", StatusInternalServerError, response.Status)
	}

	if response.Description != "Internal server error" {
		t.Errorf("Expected description 'Internal server error', got '%s'", response.Description)
	}

	customResponse := NewInternalServerError("Custom internal server error")
	if customResponse.Description != "Custom internal server error" {
		t.Errorf("Expected custom description 'Custom internal server error', got '%s'", customResponse.Description)
	}
}

func TestNewServiceUnavailableError(t *testing.T) {
	response := NewServiceUnavailableError()

	if response.Status != StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", StatusServiceUnavailable, response.Status)
	}

	if response.Description != "Service unavailable" {
		t.Errorf("Expected description 'Service unavailable', got '%s'", response.Description)
	}

	customResponse := NewServiceUnavailableError("Custom service unavailable error")
	if customResponse.Description != "Custom service unavailable error" {
		t.Errorf("Expected custom description 'Custom service unavailable error', got '%s'", customResponse.Description)
	}
}

func TestNewConflictError(t *testing.T) {
	response := NewConflictError()

	if response.Status != StatusConflict {
		t.Errorf("Expected status %d, got %d", StatusConflict, response.Status)
	}

	if response.Description != "Conflict" {
		t.Errorf("Expected description 'Conflict', got '%s'", response.Description)
	}

	customResponse := NewConflictError("Custom conflict error")
	if customResponse.Description != "Custom conflict error" {
		t.Errorf("Expected custom description 'Custom conflict error', got '%s'", customResponse.Description)
	}
}

func TestNewBadRequestError(t *testing.T) {
	response := NewBadRequestError()

	if response.Status != StatusBadRequest {
		t.Errorf("Expected status %d, got %d", StatusBadRequest, response.Status)
	}

	if response.Description != "Bad request" {
		t.Errorf("Expected description 'Bad request', got '%s'", response.Description)
	}

	customResponse := NewBadRequestError("Custom bad request error")
	if customResponse.Description != "Custom bad request error" {
		t.Errorf("Expected custom description 'Custom bad request error', got '%s'", customResponse.Description)
	}
}

func TestErrorDescriptions(t *testing.T) {
	tests := []struct {
		name     string
		response HTTPError
		expected string
	}{
		{
			name: "Bad Request Description",
			response: HTTPError{
				Status: StatusBadRequest,
			},
			expected: "Bad Request",
		},
		{
			name: "Unauthorized Description",
			response: HTTPError{
				Status: StatusUnauthorized,
			},
			expected: "Unauthorized",
		},
		{
			name: "Forbidden Description",
			response: HTTPError{
				Status: StatusForbidden,
			},
			expected: "Forbidden",
		},
		{
			name: "Not Found Description",
			response: HTTPError{
				Status: StatusNotFound,
			},
			expected: "Not Found",
		},
		{
			name: "Conflict Description",
			response: HTTPError{
				Status: StatusConflict,
			},
			expected: "Conflict",
		},
		{
			name: "Internal Server Error Description",
			response: HTTPError{
				Status: StatusInternalServerError,
			},
			expected: "Internal Server Error",
		},
		{
			name: "Service Unavailable Description",
			response: HTTPError{
				Status: StatusServiceUnavailable,
			},
			expected: "Service Unavailable",
		},
		{
			name: "Unknown Error Description",
			response: HTTPError{
				Status: 999,
			},
			expected: "Unknown Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.Descriptions(); got != tt.expected {
				t.Errorf("ErrorResponse.Descriptions() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorTypes(t *testing.T) {
	tests := []struct {
		name        string
		error       Error
		expected    int
		expectedErr string
	}{
		{"BadRequestError", BadRequestError{Description: "Bad Request"}, StatusBadRequest, "Bad Request"},
		{"UnauthorizedError", UnauthorizedError{Description: "Unauthorized"}, StatusUnauthorized, "Unauthorized"},
		{"ForbiddenError", ForbiddenError{Description: "Forbidden"}, StatusForbidden, "Forbidden"},
		{"NotFoundError", NotFoundError{Description: "Not Found"}, StatusNotFound, "Not Found"},
		{"ConflictError", ConflictError{Description: "Conflict"}, StatusConflict, "Conflict"},
		{"InternalServerError", InternalServerError{Description: "Internal Server Error"}, StatusInternalServerError, "Internal Server Error"},
		{"ServiceUnavailableError", ServiceUnavailableError{Description: "Service Unavailable"}, StatusServiceUnavailable, "Service Unavailable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.error.StatusCode(); got != tt.expected {
				t.Errorf("%s.StatusCode() = %v, want %v", tt.name, got, tt.expected)
			}

			if got := tt.error.Error(); got != tt.expectedErr {
				t.Errorf("%s.Error() = %v, want %v", tt.name, got, tt.expectedErr)
			}
		})
	}
}

func TestDefaultErrorResponses(t *testing.T) {
	if len(DefaultErrorResponses) != 2 {
		t.Errorf("Expected 2 default error responses, got %d", len(DefaultErrorResponses))
	}

	if DefaultErrorResponses[StatusBadRequest].Description != "Bad Request" {
		t.Errorf("Expected 'Bad Request' description for StatusBadRequest, got '%s'", DefaultErrorResponses[StatusBadRequest].Description)
	}

	if DefaultErrorResponses[StatusInternalServerError].Description != "Internal Server Error" {
		t.Errorf("Expected 'Internal Server Error' description for StatusInternalServerError, got '%s'", DefaultErrorResponses[StatusInternalServerError].Description)
	}
}

func TestDefaultErrorContentTypeResponses(t *testing.T) {
	expected := []string{
		"application/json",
		"application/xml",
		"multipart/form-data",
	}

	if len(DefaultErrorContentTypeResponses) != len(expected) {
		t.Errorf("Expected %d default content types, got %d", len(expected), len(DefaultErrorContentTypeResponses))
	}

	for i, contentType := range DefaultErrorContentTypeResponses {
		if contentType != expected[i] {
			t.Errorf("Expected content type '%s', got '%s'", expected[i], contentType)
		}
	}
}
