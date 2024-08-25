package lite

import (
	"strconv"

	"github.com/go-lite/lite/errors"

	"github.com/getkin/kin-openapi/openapi3"
)

type Route[T, B any] struct {
	operation   *openapi3.Operation
	path        string
	method      string
	contentType string
	statusCode  int
}

func (r Route[ResponseBody, Request]) Description(description string) Route[ResponseBody, Request] {
	r.operation.Description = description

	return r
}

func (r Route[ResponseBody, Request]) Summary(summary string) Route[ResponseBody, Request] {
	r.operation.Summary = summary

	return r
}

func (r Route[ResponseBody, Request]) OperationID(operationID string) Route[ResponseBody, Request] {
	r.operation.OperationID = operationID

	return r
}

func (r Route[ResponseBody, Request]) Deprecated() Route[ResponseBody, Request] {
	r.operation.Deprecated = true

	return r
}

func (r Route[ResponseBody, Request]) AddTags(tags ...string) Route[ResponseBody, Request] {
	r.operation.Tags = tags

	return r
}

func (r Route[ResponseBody, Request]) SetResponseContentType(contentType ContentType) Route[ResponseBody, Request] {
	r.operation.Responses.Value(strconv.Itoa(r.statusCode)).Value.Content[string(contentType)] = r.operation.Responses.
		Value(strconv.Itoa(r.statusCode)).Value.Content[r.contentType]

	delete(r.operation.Responses.Value(strconv.Itoa(r.statusCode)).Value.Content, r.contentType)

	return r
}

func (r Route[ResponseBody, Request]) AddErrorResponse(statusCode int, contentType ...ContentType) Route[ResponseBody, Request] {
	if len(contentType) == 0 {
		contentType = []ContentType{ContentType(r.contentType)}
	}

	for _, c := range contentType {
		httpError := errors.NewError(statusCode)
		description := httpError.Description()

		response := openapi3.NewResponse().WithDescription(description)

		content := openapi3.NewContentWithSchemaRef(
			openapi3.NewSchemaRef(
				"#/components/schemas/httpGenericError",
				&openapi3.Schema{}),
			[]string{string(c)},
		)
		response.WithContent(content)

		r.operation.AddResponse(statusCode, response)
	}

	return r
}
