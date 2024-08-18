package errors

import (
	"net/http"
	"testing"
)

func TestNewErrorResponse(t *testing.T) {
	id := "test-id"
	status := http.StatusBadRequest
	message := "test message"
	err := newErrorResponse(id, status, message)

	if err.ID != id {
		t.Errorf("expected %v, got %v", id, err.ID)
	}
	if err.Status != status {
		t.Errorf("expected %v, got %v", status, err.Status)
	}
	if err.Message != message {
		t.Errorf("expected %v, got %v", message, err.Message)
	}
}

func TestHTTPError_Error(t *testing.T) {
	message := "test error message"
	err := newErrorResponse("test-id", http.StatusBadRequest, message)

	if err.Error() != message {
		t.Errorf("expected %v, got %v", message, err.Error())
	}
}

func TestHTTPError_StatusCode(t *testing.T) {
	status := http.StatusBadRequest
	err := newErrorResponse("test-id", status, "test message")

	if err.StatusCode() != status {
		t.Errorf("expected %v, got %v", status, err.StatusCode())
	}
}

func TestHTTPError_Description(t *testing.T) {
	err := newErrorResponse("test-id", http.StatusBadRequest, "test message")
	expectedDescription := "Bad Request"

	if err.Description() != expectedDescription {
		t.Errorf("expected %v, got %v", expectedDescription, err.Description())
	}

	err = newErrorResponse("test-id", http.StatusInternalServerError, "test message")
	expectedDescription = "Internal Server Error"

	if err.Description() != expectedDescription {
		t.Errorf("expected %v, got %v", expectedDescription, err.Description())
	}

	err = newErrorResponse("test-id", http.StatusUnauthorized, "test message")
	expectedDescription = "Unauthorized"

	if err.Description() != expectedDescription {
		t.Errorf("expected %v, got %v", expectedDescription, err.Description())
	}

	err = newErrorResponse("test-id", http.StatusNotFound, "test message")
	expectedDescription = "Not Found"

	if err.Description() != expectedDescription {
		t.Errorf("expected %v, got %v", expectedDescription, err.Description())
	}

	err = newErrorResponse("test-id", http.StatusConflict, "test message")
	expectedDescription = "Conflict"

	if err.Description() != expectedDescription {
		t.Errorf("expected %v, got %v", expectedDescription, err.Description())
	}

	err = newErrorResponse("test-id", http.StatusForbidden, "test message")
	expectedDescription = "Forbidden"

	if err.Description() != expectedDescription {
		t.Errorf("expected %v, got %v", expectedDescription, err.Description())
	}

	err = newErrorResponse("test-id", http.StatusServiceUnavailable, "test message")
	expectedDescription = "Service Unavailable"

	if err.Description() != expectedDescription {
		t.Errorf("expected %v, got %v", expectedDescription, err.Description())
	}

	err = newErrorResponse("test-id", 999, "test message")
	expectedDescription = "Unknown Error"

	if err.Description() != expectedDescription {
		t.Errorf("expected %v, got %v", expectedDescription, err.Description())
	}
}

func TestHTTPError_SetMessage(t *testing.T) {
	err := newErrorResponse("test-id", http.StatusBadRequest, "test message")
	newMessage := "new message"
	updatedErr := err.SetMessage(newMessage)

	if updatedErr.Message != newMessage {
		t.Errorf("expected %v, got %v", newMessage, updatedErr.Message)
	}
}

func TestNewInternalServerError(t *testing.T) {
	message := "test internal server error"
	err := NewInternalServerError(message)

	if err.Status != http.StatusInternalServerError {
		t.Errorf("expected %v, got %v", http.StatusInternalServerError, err.Status)
	}

	if err.Message != message {
		t.Errorf("expected %v, got %v", message, err.Message)
	}

	err = NewInternalServerError()
	if err.Status != http.StatusInternalServerError {
		t.Errorf("expected %v, got %v", http.StatusInternalServerError, err.Status)
	}
}

func TestNewServiceUnavailableError(t *testing.T) {
	message := "test service unavailable"
	err := NewServiceUnavailableError(message)

	if err.Status != http.StatusServiceUnavailable {
		t.Errorf("expected %v, got %v", http.StatusServiceUnavailable, err.Status)
	}

	if err.Message != message {
		t.Errorf("expected %v, got %v", message, err.Message)
	}

	err = NewServiceUnavailableError()
	if err.Status != http.StatusServiceUnavailable {
		t.Errorf("expected %v, got %v", http.StatusServiceUnavailable, err.Status)
	}
}

func TestNewForbiddenError(t *testing.T) {
	message := "test forbidden"
	err := NewForbiddenError(message)

	if err.Status != http.StatusForbidden {
		t.Errorf("expected %v, got %v", http.StatusForbidden, err.Status)
	}

	if err.Message != message {
		t.Errorf("expected %v, got %v", message, err.Message)
	}

	err = NewForbiddenError()
	if err.Status != http.StatusForbidden {
		t.Errorf("expected %v, got %v", http.StatusForbidden, err.Status)
	}
}

func TestNewBadRequestError(t *testing.T) {
	message := "test bad request"
	err := NewBadRequestError(message)

	if err.Status != http.StatusBadRequest {
		t.Errorf("expected %v, got %v", http.StatusBadRequest, err.Status)
	}

	if err.Message != message {
		t.Errorf("expected %v, got %v", message, err.Message)
	}

	err = NewBadRequestError()
	if err.Status != http.StatusBadRequest {
		t.Errorf("expected %v, got %v", http.StatusBadRequest, err.Status)
	}
}

func TestNewNotFoundError(t *testing.T) {
	message := "test not found"
	err := NewNotFoundError(message)

	if err.Status != http.StatusNotFound {
		t.Errorf("expected %v, got %v", http.StatusNotFound, err.Status)
	}

	if err.Message != message {
		t.Errorf("expected %v, got %v", message, err.Message)
	}

	err = NewNotFoundError()
	if err.Status != http.StatusNotFound {
		t.Errorf("expected %v, got %v", http.StatusNotFound, err.Status)
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	message := "test unauthorized"
	err := NewUnauthorizedError(message)

	if err.Status != http.StatusUnauthorized {
		t.Errorf("expected %v, got %v", http.StatusUnauthorized, err.Status)
	}

	if err.Message != message {
		t.Errorf("expected %v, got %v", message, err.Message)
	}

	err = NewUnauthorizedError()
	if err.Status != http.StatusUnauthorized {
		t.Errorf("expected %v, got %v", http.StatusUnauthorized, err.Status)
	}
}

func TestNewConflictError(t *testing.T) {
	message := "test conflict"
	err := NewConflictError(message)

	if err.Status != http.StatusConflict {
		t.Errorf("expected %v, got %v", http.StatusConflict, err.Status)
	}

	if err.Message != message {
		t.Errorf("expected %v, got %v", message, err.Message)
	}

	err = NewConflictError()
	if err.Status != http.StatusConflict {
		t.Errorf("expected %v, got %v", http.StatusConflict, err.Status)
	}
}

func TestNewError(t *testing.T) {
	status := http.StatusForbidden
	message := "test forbidden"
	err := NewError(status, message)

	if err.Status != status {
		t.Errorf("expected %v, got %v", status, err.Status)
	}

	if err.Message != message {
		t.Errorf("expected %v, got %v", message, err.Message)
	}

	err = NewError(status)
	if err.Status != status {
		t.Errorf("expected %v, got %v", status, err.Status)
	}
}
