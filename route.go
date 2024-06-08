package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIParam struct {
	Required bool
	Example  string
	Type     string // "query", "header", "cookie"
}

type Route[T, B any] struct {
	operation   *openapi3.Operation
	path        string
	name        string
	tags        []string
	method      string
	contentType string
}

func (r Route[ResponseBody, RequestBody]) Description(description string) Route[ResponseBody, RequestBody] {
	r.operation.Description = description
	return r
}

func (r Route[ResponseBody, RequestBody]) Summary(summary string) Route[ResponseBody, RequestBody] {
	r.operation.Summary = summary
	return r
}

func (r Route[ResponseBody, RequestBody]) OperationID(operationID string) Route[ResponseBody, RequestBody] {
	r.operation.OperationID = operationID
	return r
}

func (r Route[ResponseBody, RequestBody]) Deprecated() Route[ResponseBody, RequestBody] {
	r.operation.Deprecated = true

	return r
}

func (r Route[ResponseBody, RequestBody]) AddTags(tags ...string) Route[ResponseBody, RequestBody] {
	r.operation.Tags = tags

	return r
}
