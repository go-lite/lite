package lite

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type OpenAPIParam struct {
	Required bool
	Example  string
	Type     string // "query", "header", "cookie"
}

type Route[T, B any] struct {
	Operation   *openapi3.Operation
	path        string
	name        string
	tags        []string
	method      string
	contentType string
}

func (r Route[ResponseBody, RequestBody]) Description(description string) Route[ResponseBody, RequestBody] {
	r.Operation.Description = description
	return r
}

func (r Route[ResponseBody, RequestBody]) Summary(summary string) Route[ResponseBody, RequestBody] {
	r.Operation.Summary = summary
	return r
}

func (r Route[ResponseBody, RequestBody]) OperationID(operationID string) Route[ResponseBody, RequestBody] {
	r.Operation.OperationID = operationID
	return r
}

func (r Route[ResponseBody, RequestBody]) Deprecated(b bool) Route[ResponseBody, RequestBody] {
	r.Operation.Deprecated = b

	return r
}

func (r Route[ResponseBody, RequestBody]) AddTags(tags ...string) Route[ResponseBody, RequestBody] {
	r.Operation.Tags = tags

	return r
}
