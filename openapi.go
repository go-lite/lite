package lite

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"reflect"
	"regexp"
	"slices"
	"strings"

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
	fieldType := reflect.TypeOf(*new(ResponseBody))

	err = setResponseSchema(s, operation, tag, resContentType, statusCode, fieldType)
	if err != nil {
		return nil, err
	}

	// Add error responses
	responses, _ := s.createDefaultErrorResponses()

	for code, resp := range responses {
		operation.AddResponse(code, resp)
	}

	// Remove default response
	operation.Responses.Delete("default")

	s.openAPISpec.AddOperation(routePath, method, operation)

	return operation, nil
}

func getRequiredValue(contentType string, fieldType reflect.Type, schema *openapi3.Schema) bool {
	switch fieldType.Kind() {
	case reflect.Struct:
		// if fieldType is time.Time, skip it
		if fieldType.Name() == "Time" {
			return true
		}

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
				if !slices.Contains(schema.Required, fieldName) {
					schema.Required = append(schema.Required, fieldName)
				}
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
		return true
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

		ftype := fieldType
		if fieldType.Kind() == reflect.Ptr {
			ftype = fieldType.Elem()
		}

		tagMap := parseTag(tag)

		var parameter *openapi3.Parameter
		var scheme, tpe, name string

		if pathKey, ok := tagMap["path"]; ok {
			parameter = openapi3.NewPathParameter(pathKey)

			err := setParamSchema(s, operation, pathKey, ftype.Name(), parameter, isRequired, fieldType)
			if err != nil {
				return err
			}
		} else if queryKey, ok := tagMap["query"]; ok {
			parameter = openapi3.NewQueryParameter(queryKey)

			err := setParamSchema(s, operation, queryKey, ftype.Name(), parameter, isRequired, fieldType)
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
				err := setHeaderScheme(s, operation, headerKey, ftype.Name(), parameter, fieldType, isRequired)
				if err != nil {
					return err
				}
			}
		} else if cookieKey, ok := tagMap["cookie"]; ok {
			parameter = openapi3.NewCookieParameter(cookieKey)
			err := setParamSchema(s, operation, cookieKey, ftype.Name(), parameter, isRequired, fieldType)
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
			return fmt.Errorf("unknown parameter type %s", tag)
		}
	}

	return nil
}

func setResponseSchema(
	s *App,
	operation *openapi3.Operation,
	tag string,
	resContentType string,
	statusCode int,
	fieldType reflect.Type,
) (err error) {
	responseSchema, ok := s.openAPISpec.Components.Schemas[tag]
	if !ok {
		if fieldType != any(nil) {
			responseSchema, err = generatorNewSchemaRefForValue(reflect.New(fieldType).Elem().Interface(), s.openAPISpec.Components.Schemas)
		} else {
			responseSchema, err = generatorNewSchemaRefForValue(new(any), s.openAPISpec.Components.Schemas)
		}

		if err != nil {
			return
		}

		if tag != "unknown" {
			getRequiredValue(resContentType, fieldType, responseSchema.Value)
		}

		s.openAPISpec.Components.Schemas[tag] = responseSchema
	} else {
		var newSchema *openapi3.SchemaRef

		newSchema, err = generatorNewSchemaRefForValue(reflect.New(fieldType).Elem().Interface(), s.openAPISpec.Components.Schemas)
		if err != nil {
			return
		}

		newSchemaContent := fmt.Sprintf("%v", newSchema.Value)

		if !reflect.DeepEqual(newSchema.Value, responseSchema.Value) {
			hash := computeHash(newSchemaContent)
			hashedTag := fmt.Sprintf("%s%s", tag, hash)
			s.openAPISpec.Components.Schemas[hashedTag] = newSchema

			tag = hashedTag
		}
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

	return
}

func setBodySchema(
	s *App,
	operation *openapi3.Operation,
	kind reflect.Kind,
	fieldType reflect.Type,
	fieldName string,
	contentType string,
) error {
	existingSchema, exists := s.openAPISpec.Components.Schemas[fieldName]
	if !exists {
		var err error

		if fieldType == reflect.TypeOf(multipart.FileHeader{}) || fieldType == reflect.TypeOf(&multipart.FileHeader{}) {
			fieldType = reflect.TypeOf([]byte{})
		}

		if kind == reflect.Struct {
			fieldType = updateFileHeaderFieldType(fieldType)
		}

		tp := reflect.New(fieldType).Elem().Interface()

		bodySchema, err := generatorNewSchemaRefForValue(tp, s.openAPISpec.Components.Schemas)
		if err != nil {
			return err
		}

		getRequiredValue(contentType, fieldType, bodySchema.Value)

		s.openAPISpec.Components.Schemas[fieldName] = bodySchema
	} else {
		newInstance := reflect.New(fieldType).Elem().Interface()
		newBodySchema, err := generatorNewSchemaRefForValue(newInstance, s.openAPISpec.Components.Schemas)
		if err != nil {
			return err
		}

		newSchemaContent := fmt.Sprintf("%v", newBodySchema.Value)

		if !reflect.DeepEqual(newBodySchema.Value, existingSchema.Value) {
			hash := computeHash(newSchemaContent)
			hashedFieldName := fmt.Sprintf("%s%s", fieldName, hash)
			s.openAPISpec.Components.Schemas[hashedFieldName] = newBodySchema

			fieldName = hashedFieldName
		}
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

func setHeaderScheme(
	s *App,
	operation *openapi3.Operation,
	headerName string,
	tag string,
	parameter *openapi3.Parameter,
	fieldType reflect.Type,
	isRequired bool,
) error {
	existingSchema, exists := s.openAPISpec.Components.Schemas[tag]
	if !exists {
		var err error
		newInstance := reflect.New(fieldType).Elem().Interface()

		headerSchema, err := generatorNewSchemaRefForValue(newInstance, s.openAPISpec.Components.Schemas)
		if err != nil {
			return err
		}

		headerSchema.Value.Nullable = !isRequired

		s.openAPISpec.Components.Schemas[tag] = headerSchema
	} else {
		newInstance := reflect.New(fieldType).Elem().Interface()

		headerSchema, err := generatorNewSchemaRefForValue(newInstance, s.openAPISpec.Components.Schemas)
		if err != nil {
			return err
		}

		headerSchema.Value.Nullable = !isRequired

		newSchemaContent := fmt.Sprintf("%v", headerSchema.Value)

		if !reflect.DeepEqual(headerSchema.Value, existingSchema.Value) {
			hash := computeHash(newSchemaContent)
			hashedTag := fmt.Sprintf("%s%s", tag, hash)
			s.openAPISpec.Components.Schemas[hashedTag] = headerSchema

			tag = hashedTag
		}
	}

	parameter.Schema = openapi3.NewSchemaRef(fmt.Sprintf("#/components/schemas/%s", tag), &openapi3.Schema{})
	parameter.Required = isRequired

	s.openAPISpec.Components.Parameters[headerName] = &openapi3.ParameterRef{
		Value: parameter,
	}

	operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{
		Ref: "#/components/parameters/" + headerName,
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

	s.openAPISpec.Components.SecuritySchemes[name] = securitySchemes[name]
}

// computeHash generates a SHA256 hash for the given input and returns the first 4 characters.
func computeHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return strings.ToUpper(hex.EncodeToString(hash[:])[:4])
}

func setParamSchema(
	s *App,
	operation *openapi3.Operation,
	parameterName string,
	tag string,
	parameter *openapi3.Parameter,
	isRequired bool,
	fieldType reflect.Type,
) error {
	existingSchema, exists := s.openAPISpec.Components.Schemas[tag]
	if !exists {
		var err error

		newInstance := reflect.New(fieldType).Elem().Interface()

		paramSchema, err := generatorNewSchemaRefForValue(newInstance, s.openAPISpec.Components.Schemas)
		if err != nil {
			return err
		}

		paramSchema.Value.Nullable = !isRequired

		s.openAPISpec.Components.Schemas[tag] = paramSchema
	} else {
		newInstance := reflect.New(fieldType).Elem().Interface()

		paramSchema, err := generatorNewSchemaRefForValue(newInstance, s.openAPISpec.Components.Schemas)
		if err != nil {
			return err
		}

		paramSchema.Value.Nullable = !isRequired

		newSchemaContent := fmt.Sprintf("%v", paramSchema.Value)

		if !reflect.DeepEqual(paramSchema.Value, existingSchema.Value) {
			hash := computeHash(newSchemaContent)
			hashedTag := fmt.Sprintf("%s%s", tag, hash)
			s.openAPISpec.Components.Schemas[hashedTag] = paramSchema

			tag = hashedTag
		}
	}

	parameter.Schema = openapi3.NewSchemaRef(fmt.Sprintf("#/components/schemas/%s", tag), &openapi3.Schema{})
	parameter.Required = isRequired

	s.openAPISpec.Components.Parameters[parameterName] = &openapi3.ParameterRef{
		Value: parameter,
	}

	operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{
		Ref: "#/components/parameters/" + parameterName,
	})

	return nil
}

func tagFromType(v any) string {
	if v == nil {
		return "unknown"
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
		return extractBaseName(t)
	}
}

func extractBaseName(t reflect.Type) string {
	typeName := t.String()

	// Utiliser une regex pour enlever les informations de package et les crochets
	re := regexp.MustCompile(`(?:\w+\.)?(\w+)\[([\w./-]+)\]`)

	matches := re.FindStringSubmatch(typeName)
	if len(matches) == 3 {
		// matches[1] est le nom du type générique, matches[2] est le nom du type paramétré
		genericType := matches[1]
		paramType := matches[2]
		// Extraire le nom de base du paramType en retirant le package
		if strings.Contains(paramType, "/") || strings.Contains(paramType, ".") {
			paramType = paramType[strings.LastIndex(paramType, "/")+1:]
			paramType = paramType[strings.LastIndex(paramType, ".")+1:]
		}

		return genericType + paramType
	}

	// Supprimer le préfixe du package pour les types non génériques
	if strings.Contains(typeName, ".") {
		typeName = typeName[strings.LastIndex(typeName, ".")+1:]
	}

	return typeName
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
