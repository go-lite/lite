package lite

import (
	"fmt"
	"github.com/disco07/lite-fiber/codec"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"reflect"
	"strings"
)

var generator = openapi3gen.NewGenerator(
	openapi3gen.UseAllExportedFields(),
)

func RegisterOpenAPIOperation[ResponseBody, RequestBody any](
	s *App,
	method, path, resContentType string,
	statusCode int,
) (*openapi3.Operation, error) {
	operation := &openapi3.Operation{}
	operation.OperationID = method + path

	var reqBody RequestBody
	valGen := reflect.ValueOf(&reqBody).Elem()

	countPathParams := 0

	for i := 0; i < valGen.NumField(); i++ {
		fieldGeneric := valGen.Field(i)
		fieldGenericType := fieldGeneric.Interface().(codec.TypeOf).TypeOf()

		for j := 0; j < fieldGenericType.NumField(); j++ {
			fieldTags := fieldGenericType.Field(j).Tag
			fieldType := fieldGenericType.Field(j).Type
			isPointer := fieldType.Kind() == reflect.Ptr

			// Vérifier si le champ implémente ParamType
			if paramType, ok := fieldGeneric.Interface().(codec.ParamType); ok {
				var parameter *openapi3.Parameter

				var tag, schemeTag, tpeTag, nameTag string

				switch paramType.ParamType() {
				case "query":
					tag = fieldTags.Get("query")
					parameter = openapi3.NewQueryParameter(tag)
				case "header":
					tag = fieldTags.Get("header")
					schemeTag = fieldTags.Get("scheme")
					tpeTag = fieldTags.Get("type")
					nameTag = fieldTags.Get("name")

					parameter = openapi3.NewHeaderParameter(tag)
				case "path":
					tag = fieldTags.Get("params")
					parameter = openapi3.NewPathParameter(tag)
					countPathParams++
				case "cookie":
					parameter = openapi3.NewCookieParameter(tag)
				default:
					return nil, fmt.Errorf("unknown parameter type: %s", paramType.ParamType())
				}

				ref := fmt.Sprintf("#/components/schemas/%s", tag)

				var typeScheme []string
				typeScheme = append(typeScheme, fieldType.String())

				parameter.Schema = openapi3.NewSchemaRef(ref, &openapi3.Schema{})
				parameter.Required = !isPointer

				if paramType.ParamType() == "header" {
					if strings.EqualFold(tag, "Authorization") {
						scheme := "bearer"
						tpe := "http"
						name := "Authorization"

						if nameTag != "" {
							name = nameTag
						}

						if schemeTag != "" {
							scheme = schemeTag
						}

						if tpeTag != "" {
							tpe = tpeTag
						}

						sec := openapi3.NewSecurityRequirement()
						sec[name] = []string{}

						securityScheme := openapi3.NewSecurityScheme()
						securityScheme.Type = tpe
						securityScheme.Scheme = scheme

						if operation.Security == nil {
							operation.Security = openapi3.NewSecurityRequirements()
						}

						operation.Security.With(
							sec,
						)

						var securitySchemes = make(map[string]*openapi3.SecuritySchemeRef)
						securitySchemes[name] = &openapi3.SecuritySchemeRef{
							Value: securityScheme,
						}

						s.OpenApiSpec.Components.SecuritySchemes[name] = securitySchemes[name]
					} else {
						paramSchema, ok := s.OpenApiSpec.Components.Schemas[tag]
						if !ok {
							var err error
							paramSchema, err = generator.NewSchemaRefForValue(new(string), s.OpenApiSpec.Components.Schemas)
							if err != nil {
								return operation, err
							}
							s.OpenApiSpec.Components.Schemas[tag] = paramSchema
						}

						s.OpenApiSpec.Components.Headers[tag] = &openapi3.HeaderRef{
							Value: &openapi3.Header{
								Parameter: *parameter,
							},
						}

						operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{
							Ref: "#/components/parameters/" + tag,
						})
					}
				} else {
					paramSchema, ok := s.OpenApiSpec.Components.Schemas[tag]
					if !ok {
						var err error
						paramSchema, err = generator.NewSchemaRefForValue(new(string), s.OpenApiSpec.Components.Schemas)
						if err != nil {
							return operation, err
						}
						s.OpenApiSpec.Components.Schemas[tag] = paramSchema
					}

					s.OpenApiSpec.Components.Parameters[tag] = &openapi3.ParameterRef{
						Value: parameter,
					}

					operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{
						Ref: "#/components/parameters/" + tag,
					})
				}
			}

			// Vérifier si le champ implémente ContentType
			if contentType, ok := fieldGeneric.Interface().(codec.ContentType); ok {
				bodySchema, ok := s.OpenApiSpec.Components.Schemas[fieldGenericType.Name()]
				if !ok {
					var err error
					tp := reflect.New(fieldGenericType).Elem().Interface()

					bodySchema, err = generator.NewSchemaRefForValue(tp, s.OpenApiSpec.Components.Schemas)
					if err != nil {
						return operation, err
					}

					for k := 0; k < fieldGenericType.NumField(); k++ {
						field := fieldGenericType.Field(k)
						if field.Type.Kind() != reflect.Ptr {
							fieldTag := field.Tag.Get(contentType.StructTag())
							bodySchema.Value.Required = append(bodySchema.Value.Required, fieldTag)
						}
					}

					s.OpenApiSpec.Components.Schemas[fieldGenericType.Name()] = bodySchema
				}

				requestBody := openapi3.NewRequestBody()
				content := openapi3.NewContentWithSchemaRef(
					openapi3.NewSchemaRef(fmt.Sprintf(
						"#/components/schemas/%s",
						fieldGenericType.Name(),
					), &openapi3.Schema{}),
					[]string{contentType.ContentType()},
				)

				requestBody.WithContent(content)

				operation.RequestBody = &openapi3.RequestBodyRef{
					Value: requestBody,
				}

				continue
			}
		}
	}

	routePath, pathParams := parseRoutePath(path)
	if len(pathParams) != countPathParams {
		return nil, fmt.Errorf("path params mismatch: %d declared, %d required", len(pathParams), countPathParams)
	}

	tag := tagFromType(*new(ResponseBody))

	responseSchema, ok := s.OpenApiSpec.Components.Schemas[tag]
	if !ok {
		var err error
		responseSchema, err = generator.NewSchemaRefForValue(new(ResponseBody), s.OpenApiSpec.Components.Schemas)
		if err != nil {
			return operation, err
		}
		s.OpenApiSpec.Components.Schemas[tag] = responseSchema
	}

	response := openapi3.NewResponse().WithDescription("OK")

	if responseSchema != nil {
		content := openapi3.NewContentWithSchemaRef(
			openapi3.NewSchemaRef(fmt.Sprintf(
				"#/components/schemas/%s",
				tag,
			), &openapi3.Schema{}),
			[]string{resContentType},
		)
		response.WithContent(content)
	}

	operation.AddResponse(statusCode, response)

	// Add error responses
	responses, err := s.createDefaultErrorResponses()
	if err != nil {
		return nil, err
	}

	for code, resp := range responses {
		operation.AddResponse(code, resp)
	}

	// Remove default response
	operation.Responses.Delete("default")

	s.OpenApiSpec.AddOperation(routePath, method, operation)

	return operation, nil
}

func tagFromType(v any) string {
	if v == nil {
		return "unknown-interface"
	}

	return dive(reflect.TypeOf(v), 4)
}

func dive(t reflect.Type, maxDepth int) string {
	switch t.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		if maxDepth == 0 {
			return "default"
		}
		return dive(t.Elem(), maxDepth-1)
	default:
		return t.Name()
	}
}
