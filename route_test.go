package lite

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
)

type ResponseBody struct{}

type Request struct{}

func TestRoute_Description(t *testing.T) {
	operation := &openapi3.Operation{}
	route := Route[ResponseBody, Request]{
		operation: operation,
	}

	updatedRoute := route.Description("New description")

	assert.Equal(t, "New description", updatedRoute.operation.Description)
}

func TestRoute_Summary(t *testing.T) {
	operation := &openapi3.Operation{}
	route := Route[ResponseBody, Request]{
		operation: operation,
	}

	updatedRoute := route.Summary("New summary")

	assert.Equal(t, "New summary", updatedRoute.operation.Summary)
}

func TestRoute_OperationID(t *testing.T) {
	operation := &openapi3.Operation{}
	route := Route[ResponseBody, Request]{
		operation: operation,
	}

	updatedRoute := route.OperationID("newOperationID")

	assert.Equal(t, "newOperationID", updatedRoute.operation.OperationID)
}

func TestRoute_Deprecated(t *testing.T) {
	operation := &openapi3.Operation{}
	route := Route[ResponseBody, Request]{
		operation: operation,
	}

	updatedRoute := route.Deprecated()

	assert.True(t, updatedRoute.operation.Deprecated)
}

func TestRoute_AddTags(t *testing.T) {
	operation := &openapi3.Operation{}
	route := Route[ResponseBody, Request]{
		operation: operation,
	}

	updatedRoute := route.AddTags("tag1", "tag2")

	assert.Equal(t, []string{"tag1", "tag2"}, updatedRoute.operation.Tags)
}

func TestRoute_SetResponseContentType(t *testing.T) {
	operation := &openapi3.Operation{
		Responses: &openapi3.Responses{
			Extensions: make(map[string]interface{}),
		},
	}

	operation.Responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{},
			},
		},
	})

	route := Route[ResponseBody, Request]{
		operation:   operation,
		contentType: "application/json",
		statusCode:  200,
	}

	updatedRoute := route.SetResponseContentType("application/xml")

	oldExists := updatedRoute.operation.Responses.Value("200")
	newExists := updatedRoute.operation.Responses.Value("200")

	assert.Nil(t, oldExists.Value.Content["application/json"])
	assert.NotNil(t, newExists.Value.Content["application/xml"])
}

func TestRoute_AddErrorResponse(t *testing.T) {
	operation := &openapi3.Operation{
		Responses: &openapi3.Responses{
			Extensions: make(map[string]interface{}),
		},
	}

	route := Route[ResponseBody, Request]{
		operation:   operation,
		contentType: "application/json",
		statusCode:  200,
	}

	route = route.AddErrorResponse(400)

	assert.Equal(t, "Bad Request", *route.operation.Responses.Value("400").Value.Description)
}
