package lite

import (
	"github.com/getkin/kin-openapi/openapi3"
	"go.opentelemetry.io/otel/trace"
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

	tracerProvider trace.TracerProvider
	cfg            Config
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

// Param registers a parameter for the route.
// The paramType can be "query", "header" or "cookie".
// [Cookie], [Header], [QueryParam] are shortcuts for Param.
func (r Route[ResponseBody, RequestBody]) Param(paramType, name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	openapiParam := openapi3.NewHeaderParameter(name)
	openapiParam.Description = description
	openapiParam.Schema = openapi3.NewStringSchema().NewRef()
	openapiParam.In = paramType

	for _, param := range params {
		if param.Required {
			openapiParam.Required = param.Required
		}
		if param.Example != "" {
			openapiParam.Example = param.Example
		}
	}

	r.Operation.AddParameter(openapiParam)

	return r
}

func (r Route[ResponseBody, RequestBody]) Header(name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	r.Param(openapi3.ParameterInHeader, name, description, params...)
	return r
}

func (r Route[ResponseBody, RequestBody]) Cookie(name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	r.Param(openapi3.ParameterInCookie, name, description, params...)
	return r
}

func (r Route[ResponseBody, RequestBody]) QueryParam(name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	r.Param(openapi3.ParameterInQuery, name, description, params...)
	return r
}

func (r Route[ResponseBody, RequestBody]) PathParam(name, description string, params ...OpenAPIParam) Route[ResponseBody, RequestBody] {
	r.Param(openapi3.ParameterInPath, name, description, params...)
	return r
}

func (r Route[ResponseBody, RequestBody]) Response(status int, description string) Route[ResponseBody, RequestBody] {
	response := openapi3.NewResponse()
	response.Description = &description

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

func (r Route[ResponseBody, RequestBody]) AddServer(url, desc string) {
	var servers []*openapi3.Server

	servers = append(servers, &openapi3.Server{
		URL:         url,
		Description: desc,
	})

	r.Operation.Servers = (*openapi3.Servers)(&servers)
}
