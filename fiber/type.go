package fiber

import (
	"log"
	"reflect"
	"swagger/swagger"
)

type ContentTyper interface {
	ContentType() string
}

type AsHeader[T any] struct {
	Value T
}

type AsPathParam[T any] struct {
	Value T
}

type AsQueryParam[T any] struct {
	Value T
}

type AsJSON[T any] struct {
	Value T
}

func (a AsJSON[T]) ContentType() string {
	return "application/json"
}

type AsPDF[T any] struct {
	Value T
}

func (a AsPDF[T]) ContentType() string {
	return "application/pdf"
}

type AsMultiPart[T any] struct {
	Value T
}

func (a AsMultiPart[T]) ContentType() string {
	return "multipart/form-data"
}

type AsTextPlain[T any] struct {
	Value T
}

func (a AsTextPlain[T]) ContentType() string {
	return "text/plain"
}

func getContentType(field reflect.Value) (string, bool) {
	fieldInterface := field.Interface()
	if typer, ok := fieldInterface.(ContentTyper); ok {
		return typer.ContentType(), true
	}
	return "", false
}

func createInstanceFromType(t reflect.Type) reflect.Value {
	// Crée une nouvelle instance du type
	instance := reflect.New(t).Elem()
	return instance
}

func generateRequestSchema[T any]() (swagger.Value, []swagger.Header, []swagger.PathParameter, []swagger.QueryParameter, string) {
	tType := reflect.TypeOf((*T)(nil)).Elem()
	val := createInstanceFromType(tType)

	schema := generateSchemaFromType(tType)

	var headers []swagger.Header
	var pathParams []swagger.PathParameter
	var queryParams []swagger.QueryParameter
	var contentType string

	for i := 0; i < tType.NumField(); i++ {
		fieldTpe := tType.Field(i)
		fieldValue := val.Field(i)
		tag := fieldTpe.Tag

		var fieldSchema swagger.Value
		fieldType := fieldTpe.Type
		if fieldType.Kind() == reflect.Struct {
			if fieldType.NumField() > 0 {
				genericType := fieldType.Field(0).Type
				fieldSchema = generateSchemaFromType(genericType)
			}
		}

		if headerTag, ok := tag.Lookup("header"); ok {
			headers = append(headers, swagger.Header{
				Key:   headerTag,
				Value: fieldSchema,
			})

			continue
		}

		if paramTag, ok := tag.Lookup("param"); ok {
			pathParams = append(pathParams, swagger.PathParameter{
				Key:   paramTag,
				Value: fieldSchema,
			})

			continue
		}

		if queryTag, ok := tag.Lookup("query"); ok {
			queryParams = append(queryParams, swagger.QueryParameter{
				Key:   queryTag,
				Value: fieldSchema,
			})

			continue
		}

		schema = fieldSchema

		contentType, _ = getContentType(fieldValue)
	}

	return schema, headers, pathParams, queryParams, contentType
}

func generateSchema[T any]() (swagger.Value, string) {
	tType := reflect.TypeOf((*T)(nil)).Elem()
	contentType, _ := "application/json", false
	return generateSchemaFromType(tType), contentType
}

func generateSchemaFromType(tType reflect.Type) swagger.Value {
	schema := swagger.Value{
		Type:       swagger.ValueTypeObject,
		Properties: make(map[string]swagger.Value),
	}

	switch tType.Kind() {
	case reflect.Struct:
		for i := 0; i < tType.NumField(); i++ {
			field := tType.Field(i)
			fieldName := field.Name
			tag := field.Tag
			if headerTag, ok := tag.Lookup("reqHeader"); ok {
				log.Println("header", headerTag)
			}

			if jsonTag := tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				fieldName = jsonTag
			}
			fieldSchema := generateSchemaFromType(field.Type)
			schema.Properties[fieldName] = fieldSchema
		}
	case reflect.String:
		schema.Type = swagger.ValueTypeString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		schema.Type = swagger.ValueTypeInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		schema.Type = swagger.ValueTypeUint
	case reflect.Float32, reflect.Float64:
		schema.Type = swagger.ValueTypeFloat
	case reflect.Bool:
		schema.Type = swagger.ValueTypeBool
	case reflect.Slice:
		schema.Type = swagger.ValueTypeArray
		schema.Items = &swagger.Value{
			Type: swagger.ValueTypeAny,
		}
		if tType.Elem().Kind() == reflect.Struct {
			elemSchema := generateSchemaFromType(tType.Elem())
			schema.Items = &elemSchema
		}
	case reflect.Map:
		schema.Type = swagger.ValueTypeMap
		keySchema := generateSchemaFromType(tType.Key())
		valueSchema := generateSchemaFromType(tType.Elem())
		schema.Keys = &keySchema
		schema.Values = &valueSchema
	case reflect.Ptr:
		schema = generateSchemaFromType(tType.Elem())
		schema.Nullable = true
	default:
		schema.Type = swagger.ValueTypeAny
	}

	return schema
}
