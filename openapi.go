package lite

import (
	"fmt"
	"mime/multipart"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
)

var generator = openapi3gen.NewGenerator(
	openapi3gen.UseAllExportedFields(),
)

var generatorNewSchemaRefForValue = generator.NewSchemaRefForValue

func registerOpenAPIOperation[ResponseBody, RequestBody any](
	s *App,
	method, path, resContentType string,
	statusCode int,
) (operation *openapi3.Operation, err error) {
	operation = openapi3.NewOperation()
	operation.OperationID = method + path

	var reqBody RequestBody
	valGen := reflect.ValueOf(&reqBody).Elem()
	kind := valGen.Kind()

	switch kind {
	case reflect.Struct:
		err = register(s, operation, valGen)
		if err != nil {
			return nil, err
		}
	case reflect.Slice:
		if valGen.Type().Elem().Kind() == reflect.Uint8 {
			err = setBodySchema(s, operation, kind, valGen.Type(), valGen.Type().Elem().Name(), "application/octet-stream")
			if err != nil {
				return nil, err
			}
		}
	case reflect.String:
		err = setBodySchema(s, operation, kind, valGen.Type(), valGen.Type().Name(), "text/plain")
		if err != nil {
			return nil, err
		}
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		fallthrough
	default:
	}

	routePath, _ := parseRoutePath(path)

	tag := tagFromType(*new(ResponseBody))

	responseSchema, ok := s.OpenAPISpec.Components.Schemas[tag]
	if !ok {
		responseSchema, err = generatorNewSchemaRefForValue(new(ResponseBody), s.OpenAPISpec.Components.Schemas)
		if err != nil {
			return operation, err
		}

		fieldGenericType := reflect.TypeOf(*new(ResponseBody))

		if tag != "unknown-interface" {
			getRequiredValue(resContentType, fieldGenericType, responseSchema.Value)
		}

		s.OpenAPISpec.Components.Schemas[tag] = responseSchema
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

	s.OpenAPISpec.AddOperation(routePath, method, operation)

	return operation, nil
}

func getRequiredValue(contentType string, fieldType reflect.Type, schema *openapi3.Schema) bool {
	switch fieldType.Kind() {
	case reflect.Struct:
		for k := 0; k < fieldType.NumField(); k++ {
			field := fieldType.Field(k)
			fieldName := field.Name

			if field.Tag.Get(getStructTag(contentType)) != "" {
				if contentType != "application/json" {
					jsonFieldName := field.Tag.Get(getStructTag("application/json"))
					if jsonFieldName != "" {
						fieldName = jsonFieldName
					}

					updateKey(schema.Properties, fieldName, field.Tag.Get(getStructTag(contentType)))
				}

				fieldName = field.Tag.Get(getStructTag(contentType))
			}

			ok := getRequiredValue(contentType, field.Type, schema.Properties[fieldName].Value)
			if ok {
				schema.Required = append(schema.Required, fieldName)
			}
		}

		return true
	case reflect.Array, reflect.Slice:
		if fieldType.Elem().Kind() == reflect.Uint8 {
			return true
		}

		return getRequiredValue(contentType, fieldType.Elem(), schema.Items.Value)
	case reflect.Map:
		getRequiredValue(contentType, fieldType.Elem(), schema.AdditionalProperties.Schema.Value)
		return false
	case reflect.Interface:
		return false
	case reflect.Ptr:
		return false
	case reflect.Invalid, reflect.Uintptr, reflect.Complex64, reflect.Complex128,
		reflect.Chan, reflect.Func, reflect.UnsafePointer:
		panic("not implemented")
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
		fallthrough
	default:
		return true
	}
}

func updateKey(properties openapi3.Schemas, key string, newKey string) {
	schema := properties[key]
	properties[newKey] = schema

	if key != newKey {
		delete(properties, key)
	}
}

func register(s *App, operation *openapi3.Operation, dstVal reflect.Value) error {
	dstType := dstVal.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		fieldVal := dstVal.Field(i)
		fieldType := field.Type
		tag := field.Tag.Get("lite")
		kind := fieldVal.Kind()

		// check if kind is a pointer and elem is a not string, float, int, bool continue to next field
		switch kind {
		case reflect.Ptr:
			switch fieldVal.Elem().Kind() {
			case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map, reflect.Complex64, reflect.Complex128,
				reflect.Uintptr, reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Ptr:
				return fmt.Errorf("not implemented")
			case reflect.Invalid:
				fallthrough
			case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
				reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64,
				reflect.Interface, reflect.String:
				fallthrough
			default:
			}
		case reflect.Invalid, reflect.Uintptr, reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Complex64,
			reflect.Complex128:
			return fmt.Errorf("not implemented")
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64,
			reflect.Array, reflect.Interface, reflect.Map, reflect.Slice, reflect.String, reflect.Struct:
			fallthrough
		default:
		}

		if kind == reflect.Struct && tag == "" {
			// Recursively handle nested structs
			if err := register(s, operation, fieldVal); err != nil {
				return err
			}

			continue
		}

		isRequired := kind != reflect.Ptr

		if tag == "" {
			tag = field.Name
		}

		tagMap := parseTag(tag)

		var parameter *openapi3.Parameter
		var scheme, tpe, name string

		if pathKey, ok := tagMap["path"]; ok {
			parameter = openapi3.NewPathParameter(pathKey)

			err := setParamSchema(s, operation, pathKey, parameter, isRequired, fieldType)
			if err != nil {
				return err
			}
		} else if queryKey, ok := tagMap["query"]; ok {
			parameter = openapi3.NewQueryParameter(queryKey)
			err := setParamSchema(s, operation, queryKey, parameter, isRequired, fieldType)
			if err != nil {
				return err
			}
		} else if headerKey, ok := tagMap["header"]; ok {
			parameter = openapi3.NewHeaderParameter(headerKey)
			parameter.Required = isRequired
			var isAuth bool

			if _, isAuth = tagMap["isauth"]; isAuth {
				tpe = "http"
				name = "Authorization"
				scheme = "bearer"
			}

			if valueScheme, ok := tagMap["scheme"]; ok {
				scheme = valueScheme
			}

			if valueName, ok := tagMap["name"]; ok {
				name = valueName
			}

			if isAuth {
				setSecurityScheme(s, operation, name, tpe, scheme)
			} else {
				err := setHeaderScheme(s, operation, tag, parameter)
				if err != nil {
					return err
				}
			}
		} else if cookieKey, ok := tagMap["cookie"]; ok {
			parameter = openapi3.NewCookieParameter(cookieKey)
			err := setParamSchema(s, operation, cookieKey, parameter, isRequired, fieldType)
			if err != nil {
				return err
			}
		} else if reqKey, ok := tagMap["req"]; ok && reqKey == "body" {
			contentType := "application/json"

			if len(tagMap) > 2 {
				panic("invalid tag")
			}

			if len(tagMap) == 2 {
				for key := range tagMap {
					if key != "req" {
						contentType = key

						break
					}
				}
			}

			fieldName := field.Name

			if kind == reflect.Struct {
				fieldName = fieldType.Name()
			}

			err := setBodySchema(s, operation, kind, fieldType, fieldName, contentType)
			if err != nil {
				return err
			}

			continue
		} else {
			return fmt.Errorf("unknown parameter type")
		}
	}

	return nil
}

func setBodySchema(
	s *App,
	operation *openapi3.Operation,
	kind reflect.Kind,
	fieldType reflect.Type,
	fieldName string,
	contentType string,
) error {
	if kind != reflect.Struct && kind != reflect.String &&
		!(kind == reflect.Slice && fieldType.Elem().Kind() == reflect.Uint8) {
		return fmt.Errorf("invalid request body type %s", kind)
	}

	_, ok := s.OpenAPISpec.Components.Schemas[fieldName]
	if !ok {
		var err error

		if kind == reflect.Struct {
			fieldType = updateFileHeaderFieldType(fieldType)
		}

		tp := reflect.New(fieldType).Elem().Interface()

		bodySchema, err := generatorNewSchemaRefForValue(tp, s.OpenAPISpec.Components.Schemas)
		if err != nil {
			return err
		}

		getRequiredValue(contentType, fieldType, bodySchema.Value)

		s.OpenAPISpec.Components.Schemas[fieldName] = bodySchema
	}

	requestBody := openapi3.NewRequestBody()
	content := openapi3.NewContentWithSchemaRef(
		openapi3.NewSchemaRef(fmt.Sprintf(
			"#/components/schemas/%s",
			fieldName,
		), &openapi3.Schema{}),
		[]string{contentType},
	)

	requestBody.WithContent(content)

	operation.RequestBody = &openapi3.RequestBodyRef{
		Value: requestBody,
	}

	return nil
}

// check if reflect.Type (struct) has field type multipart.FileHeader or *multipart.FileHeader
// if true, update type to []byte.
func updateFileHeaderFieldType(fieldType reflect.Type) reflect.Type {
	var fields []reflect.StructField

	for i := 0; i < fieldType.NumField(); i++ {
		field := fieldType.Field(i)

		if field.Type == reflect.TypeOf(multipart.FileHeader{}) || field.Type == reflect.TypeOf(&multipart.FileHeader{}) {
			field.Type = reflect.TypeOf([]byte{})
		}

		fields = append(fields, field)
	}

	return reflect.StructOf(fields)
}

func setHeaderScheme(s *App, operation *openapi3.Operation, tag string, parameter *openapi3.Parameter) error {
	_, ok := s.OpenAPISpec.Components.Schemas[tag]
	if !ok {
		var err error

		headerSchema, err := generatorNewSchemaRefForValue(new(string), s.OpenAPISpec.Components.Schemas)
		if err != nil {
			return err
		}

		s.OpenAPISpec.Components.Schemas[tag] = headerSchema
	}

	s.OpenAPISpec.Components.Headers[tag] = &openapi3.HeaderRef{
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

	securitySchemes := make(map[string]*openapi3.SecuritySchemeRef)
	securitySchemes[name] = &openapi3.SecuritySchemeRef{
		Value: securityScheme,
	}

	s.OpenAPISpec.Components.SecuritySchemes[name] = securitySchemes[name]
}

func setParamSchema(
	s *App,
	operation *openapi3.Operation,
	tag string,
	parameter *openapi3.Parameter,
	isRequired bool,
	fieldType reflect.Type,
) error {
	ref := fmt.Sprintf("#/components/schemas/%s", tag)

	parameter.Schema = openapi3.NewSchemaRef(ref, &openapi3.Schema{})
	parameter.Required = isRequired

	_, ok := s.OpenAPISpec.Components.Schemas[tag]
	if !ok {
		var err error
		newInstance := reflect.New(fieldType).Elem().Interface()

		paramSchema, err := generatorNewSchemaRefForValue(newInstance, s.OpenAPISpec.Components.Schemas)
		if err != nil {
			return err
		}

		s.OpenAPISpec.Components.Schemas[tag] = paramSchema
	}

	s.OpenAPISpec.Components.Parameters[tag] = &openapi3.ParameterRef{
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
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Interface, reflect.String, reflect.Struct:
		fallthrough
	default:
		return t.Name()
	}
}

// get struct tag from content type
func getStructTag(contentType string) string {
	switch contentType {
	case "application/json":
		return "json"
	case "application/xml", "text/xml":
		return "xml"
	case "application/x-www-form-urlencoded", "multipart/form-data":
		return "form"
	case "text/plain":
		return "txt"
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
