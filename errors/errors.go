package errors

import (
	"net/http"

	"github.com/google/uuid"
)

type HTTPError struct {
	ID      string `form:"id"      json:"id"      xml:"id"`
	Status  int    `form:"status"  json:"status"  xml:"status"`
	Message string `form:"message" json:"message" xml:"message"`
}

func newErrorResponse(id string, status int, message string) HTTPError {
	return HTTPError{
		ID:      id,
		Status:  status,
		Message: message,
	}
}

func (e HTTPError) Error() string {
	return e.Message
}

func (e HTTPError) StatusCode() int {
	return e.Status
}

func (e HTTPError) Description() string {
	switch e.Status {
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusConflict:
		return "Conflict"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	case http.StatusServiceUnavailable:
		return "Service Unavailable"
	default:
		return "Unknown Error"
	}
}

func (e HTTPError) SetMessage(message string) HTTPError {
	e.Message = message

	return e
}

var DefaultErrorResponses = map[int]HTTPError{
	http.StatusBadRequest:          newErrorResponse(uuid.NewString(), http.StatusBadRequest, "Bad Request"),
	http.StatusInternalServerError: newErrorResponse(uuid.NewString(), http.StatusInternalServerError, "Internal Server Error"),
}

var DefaultErrorContentTypeResponses = []string{
	"application/json",
	"application/xml",
	"multipart/form-data",
}

func NewInternalServerError(message ...string) HTTPError {
	if len(message) > 0 {
		return newErrorResponse(uuid.NewString(), http.StatusInternalServerError, message[0])
	}

	return DefaultErrorResponses[http.StatusInternalServerError]
}

func NewBadRequestError(message ...string) HTTPError {
	if len(message) > 0 {
		return newErrorResponse(uuid.NewString(), http.StatusBadRequest, message[0])
	}

	return DefaultErrorResponses[http.StatusBadRequest]
}

func NewForbiddenError(message ...string) HTTPError {
	if len(message) > 0 {
		return newErrorResponse(uuid.NewString(), http.StatusForbidden, message[0])
	}

	return newErrorResponse(uuid.NewString(), http.StatusForbidden, "Forbidden")
}

func NewServiceUnavailableError(message ...string) HTTPError {
	if len(message) > 0 {
		return newErrorResponse(uuid.NewString(), http.StatusServiceUnavailable, message[0])
	}

	return newErrorResponse(uuid.NewString(), http.StatusServiceUnavailable, "Service Unavailable")
}

func NewNotFoundError(message ...string) HTTPError {
	if len(message) > 0 {
		return newErrorResponse(uuid.NewString(), http.StatusNotFound, message[0])
	}

	return newErrorResponse(uuid.NewString(), http.StatusNotFound, "Not Found")
}

func NewUnauthorizedError(message ...string) HTTPError {
	if len(message) > 0 {
		return newErrorResponse(uuid.NewString(), http.StatusUnauthorized, message[0])
	}

	return newErrorResponse(uuid.NewString(), http.StatusUnauthorized, "Unauthorized")
}

func NewConflictError(message ...string) HTTPError {
	if len(message) > 0 {
		return newErrorResponse(uuid.NewString(), http.StatusConflict, message[0])
	}

	return newErrorResponse(uuid.NewString(), http.StatusConflict, "Conflict")
}

func NewError(status int, message ...string) HTTPError {
	if len(message) > 0 {
		return newErrorResponse(uuid.NewString(), status, message[0])
	}

	return newErrorResponse(uuid.NewString(), status, "")
}
