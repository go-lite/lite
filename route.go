package lite

import (
	"github.com/getkin/kin-openapi/openapi3"
	"strconv"
)

type Route[T, B any] struct {
	operation   *openapi3.Operation
	path        string
	name        string
	tags        []string
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

func (r Route[ResponseBody, Request]) SetResponseContentType(contentType string) Route[ResponseBody, Request] {
	r.operation.Responses.Value(strconv.Itoa(r.statusCode)).Value.Content[contentType] = r.operation.Responses.Value("200").Value.Content[r.contentType]

	return r
}
