package openapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"reflect"
	"strings"
)

var generator = openapi3gen.NewGenerator(
	openapi3gen.UseAllExportedFields(),
)

func registerOpenAPIOperation[ResponseBody, RequestBody any](
	s *App,
	method, path, resContentType string,
	statusCode int,
) (*openapi3.Operation, error) {
	operation := &openapi3.Operation{}
	operation.OperationID = method + path

	var reqBody RequestBody
	valGen := reflect.ValueOf(&reqBody).Elem()
	kind := valGen.Kind()

	if kind != reflect.Struct {
		return nil, fmt.Errorf("request body must be a struct")
	}

	err := register(s, operation, valGen)
	if err != nil {
		return nil, err
	}

	routePath, _ := parseRoutePath(path)

	tag := tagFromType(*new(ResponseBody))

	responseSchema, ok := s.OpenApiSpec.Components.Schemas[tag]
	if !ok {
		var err error
		responseSchema, err = generator.NewSchemaRefForValue(new(ResponseBody), s.OpenApiSpec.Components.Schemas)
		if err != nil {
			return operation, err
		}

		fieldGenericType := reflect.TypeOf(*new(ResponseBody))

		if fieldGenericType.Kind() == reflect.Struct {
			for k := 0; k < fieldGenericType.NumField(); k++ {
				field := fieldGenericType.Field(k)
				if field.Type.Kind() != reflect.Ptr {
					fieldTag := field.Tag.Get(getStructTag(resContentType))
					responseSchema.Value.Required = append(responseSchema.Value.Required, fieldTag)
				}
			}
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

func register(s *App, operation *openapi3.Operation, dstVal reflect.Value) error {
	dstType := dstVal.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		fieldVal := dstVal.Field(i)
		fieldType := field.Type
		tag := field.Tag.Get("lite")

		// check if kind is a pointer and elem is a not string, float, int, bool continue to next field
		switch fieldVal.Kind() {
		case reflect.Ptr:
			switch fieldVal.Elem().Kind() {
			case reflect.Struct, reflect.Slice, reflect.Map:
				panic("not implemented")
			default:
			}

		case reflect.Invalid, reflect.Uintptr, reflect.Chan, reflect.Func, reflect.UnsafePointer:
			panic("not implemented")
		default:
		}

		if fieldVal.Kind() == reflect.Struct && tag == "" {
			// Recursively handle nested structs
			if err := register(s, operation, fieldVal); err != nil {
				return err
			}

			continue
		}

		isRequired := fieldType.Kind() != reflect.Ptr

		if tag == "" {
			tag = field.Name
		}

		tagParts := strings.Split(tag, "=")

		var parameter *openapi3.Parameter
		var scheme, tpe, name string

		tagType, tagName := tagParts[0], tagParts[1]

		ref := fmt.Sprintf("#/components/schemas/%s", tagName)

		switch tagType {
		case "path":
			parameter = openapi3.NewPathParameter(tagName)
			err := setParamSchema(s, operation, tagName, parameter, isRequired, ref)
			if err != nil {
				return err
			}
		case "query":
			parameter = openapi3.NewQueryParameter(tagName)
			err := setParamSchema(s, operation, tagName, parameter, isRequired, ref)
			if err != nil {
				return err
			}
		case "header":
			parameter = openapi3.NewHeaderParameter(tagName)
			parameter.Required = isRequired
			var isAuth bool
			if strings.Contains(tagName, ",") {
				tagDetails := strings.Split(tagName, ",")
				tagName = tagDetails[0]

				for _, detail := range tagDetails[1:] {
					switch detail {
					case "isauth":
						isAuth = true
						tpe = "http"
						name = "Authorization"
						scheme = "bearer"
					case "scheme":
						scheme = detail
					}
				}
			}

			if isAuth {
				setSecurityScheme(s, operation, name, tpe, scheme)
			} else {
				err := setHeaderScheme(s, operation, tag, parameter)
				if err != nil {
					return err
				}
			}
		case "cookie":
			parameter = openapi3.NewCookieParameter(tagName)
			err := setParamSchema(s, operation, tagName, parameter, isRequired, ref)
			if err != nil {
				return err
			}
		case "req":
			//if tagName == "body" {
			//	bodySchema, ok := s.OpenApiSpec.Components.Schemas[field.Name]
			//	if !ok {
			//		var err error
			//		tp := reflect.New(fieldGenericType).Elem().Interface()
			//
			//		bodySchema, err = generator.NewSchemaRefForValue(tp, s.OpenApiSpec.Components.Schemas)
			//		if err != nil {
			//			return operation, err
			//		}
			//
			//		for k := 0; k < fieldGenericType.NumField(); k++ {
			//			field := fieldGenericType.Field(k)
			//			if field.Type.Kind() != reflect.Ptr {
			//				fieldTag := field.Tag.Get(contentType.StructTag())
			//				bodySchema.Value.Required = append(bodySchema.Value.Required, fieldTag)
			//			}
			//		}
			//
			//		s.OpenApiSpec.Components.Schemas[fieldGenericType.Name()] = bodySchema
			//	}
			//
			//	requestBody := openapi3.NewRequestBody()
			//	content := openapi3.NewContentWithSchemaRef(
			//		openapi3.NewSchemaRef(fmt.Sprintf(
			//			"#/components/schemas/%s",
			//			fieldGenericType.Name(),
			//		), &openapi3.Schema{}),
			//		[]string{contentType.ContentType()},
			//	)
			//
			//	requestBody.WithContent(content)
			//
			//	operation.RequestBody = &openapi3.RequestBodyRef{
			//		Value: requestBody,
			//	}
			//
			//	continue
			//}

		default:
			return fmt.Errorf("unknown parameter type: %s", tagType)
		}
	}

	return nil
}

func setHeaderScheme(s *App, operation *openapi3.Operation, tag string, parameter *openapi3.Parameter) error {
	paramSchema, ok := s.OpenApiSpec.Components.Schemas[tag]
	if !ok {
		var err error
		paramSchema, err = generator.NewSchemaRefForValue(new(string), s.OpenApiSpec.Components.Schemas)
		if err != nil {
			return err
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
	return nil
}

func setSecurityScheme(s *App, operation *openapi3.Operation, name string, tpe string, scheme string) {
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
}

func setParamSchema(s *App, operation *openapi3.Operation, tag string, parameter *openapi3.Parameter, isRequired bool, ref string) error {
	parameter.Schema = openapi3.NewSchemaRef(ref, &openapi3.Schema{})
	parameter.Required = isRequired

	paramSchema, ok := s.OpenApiSpec.Components.Schemas[tag]
	if !ok {
		var err error
		paramSchema, err = generator.NewSchemaRefForValue(new(string), s.OpenApiSpec.Components.Schemas)
		if err != nil {
			return err
		}
		s.OpenApiSpec.Components.Schemas[tag] = paramSchema
	}

	s.OpenApiSpec.Components.Parameters[tag] = &openapi3.ParameterRef{
		Value: parameter,
	}

	operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{
		Ref: "#/components/parameters/" + tag,
	})

	return nil
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

// get struct tag from content type
func getStructTag(contentType string) string {
	switch contentType {
	case "application/json":
		return "json"
	case "application/xml":
		return "xml"
	case "application/x-www-form-urlencoded", "multipart/form-data":
		return "form"
	case "text/plain":
		return "text"
	case "application/octet-stream":
		return "binary"
	case "application/pdf":
		return "pdf"
	case "image/png":
		return "png"
	case "image/jpeg":
		return "jpeg"
	default:
		return "json"
	}
}
