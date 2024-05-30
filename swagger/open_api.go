package swagger

import (
	"github.com/getkin/kin-openapi/openapi3"
	"swagger/lite"
)

func RegisterOpenAPIOperation[T, B any](s *lite.App, method, path string) (*openapi3.Operation, error) {
	operation := &openapi3.Operation{}
	operation.OperationID = method + path
	operation.Tags = s.GetTags()

	// Add the operation to the OpenAPI spec
	if s.OpenApiSpec.Paths == nil {
		s.OpenApiSpec.Paths = &openapi3.Paths{}
	}

	return operation, nil
}
