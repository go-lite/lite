package lite

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type Route[T, B any] struct {
	operation   *openapi3.Operation
	path        string
	name        string
	tags        []string
	method      string
	contentType string
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
