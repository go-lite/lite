package lite

import (
	"github.com/disco07/lite-fiber/codec"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"net/http"
)

type HTTPError struct {
	ID      string `json:"id" xml:"id" form:"id"`
	Status  int    `json:"status" xml:"status" form:"status"`
	Message string `json:"message" xml:"message" form:"message"`
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
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusConflict:
		return "Conflict"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Unknown Error"
	}
}

func (e HTTPError) tag() string {
	switch e.Status {
	case http.StatusBadRequest:
		return "swaggerDescriptionBadRequest"
	case http.StatusUnauthorized:
		return "swaggerDescriptionUnauthorized"
	case http.StatusNotFound:
		return "swaggerDescriptionNotFound"
	case http.StatusConflict:
		return "swaggerDescriptionConflict"
	case http.StatusInternalServerError:
		return "swaggerDescriptionInternalServerError"
	default:
		return "swaggerDescriptionUnknownError"
	}
}

// Set error messages
func (e HTTPError) SetMessage(message string) HTTPError {
	e.Message = message

	return e
}

var defaultErrorResponses = map[int]HTTPError{
	http.StatusBadRequest:          newErrorResponse(uuid.NewString(), http.StatusBadRequest, "Bad Request"),
	http.StatusUnauthorized:        newErrorResponse(uuid.NewString(), http.StatusUnauthorized, "Unauthorized"),
	http.StatusNotFound:            newErrorResponse(uuid.NewString(), http.StatusNotFound, "Not Found"),
	http.StatusInternalServerError: newErrorResponse(uuid.NewString(), http.StatusInternalServerError, "Internal Server Error"),
}

var defaultErrorContentTypeResponses = map[string]codec.Encoder[HTTPError]{
	"application/json":    codec.AsJSON[HTTPError]{},
	"application/xml":     codec.AsXML[HTTPError]{},
	"multipart/form-data": codec.AsMultiPart[HTTPError]{},
}

func newOpenAPIErrorResponse(err HTTPError) *openapi3.Response {
	response := openapi3.NewResponse().WithDescription(err.Description())

	return response
}

func newOpenAPIErrorResponses(errs ...HTTPError) map[string]*openapi3.Response {
	responses := make(map[string]*openapi3.Response, len(errs))
	for _, err := range errs {
		responses[err.ID] = newOpenAPIErrorResponse(err)
	}

	return responses
}
