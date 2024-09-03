package lite

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

func deserializeRequests(ctx *fasthttp.RequestCtx, dst any, params map[string]string) error {
	return deserialize(ctx, reflect.ValueOf(dst).Elem(), params)
}

func deserialize(ctx *fasthttp.RequestCtx, dstVal reflect.Value, params map[string]string) error {
	dstType := dstVal.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		fieldVal := dstVal.Field(i)
		tag := field.Tag.Get("lite")

		if fieldVal.Kind() == reflect.Struct && tag == "" {
			if err := deserialize(ctx, fieldVal, params); err != nil {
				return err
			}

			continue
		}

		if tag == "" {
			return InternalServerError{
				Context:     "/api/contexts/DeserializationError",
				Type:        "DeserializationError",
				Status:      StatusInternalServerError,
				Title:       "Internal server error",
				Description: "Missing tag for field " + field.Name,
			}
		}

		tagMap := parseTag(tag)

		if val, ok := tagMap["req"]; ok && val == "body" {
			if err := deserializeBody(ctx, fieldVal); err != nil {
				return err
			}
		}

		var valueStr string

		switch {
		case tagMap["params"] != "":
			if params != nil {
				paramsKey := tagMap["params"]

				paramsValue, ok := params[paramsKey]
				if !ok {
					return InternalServerError{
						Context:     "/api/contexts/DeserializationError",
						Type:        "DeserializationError",
						Title:       "Internal server error",
						Description: "Missing params parameter: " + paramsKey,
						Violations: []Violation{
							{
								PropertyPath: paramsKey,
								Message:      "Missing params parameter: " + paramsKey,
							},
						},
					}
				}

				valueStr = paramsValue
			}
		case tagMap["query"] != "":
			queryKey := tagMap["query"]
			if value := ctx.QueryArgs().Peek(queryKey); len(value) > 0 {
				valueStr = string(value)
			}
		case tagMap["header"] != "":
			headerKey := tagMap["header"]

			if tagMap["type"] == "apiKey" {
				headerKey = tagMap["name"]
			}

			if value := ctx.Request.Header.Peek(headerKey); len(value) > 0 {
				valueStr = string(value)
			}
		}

		if valueStr != "" {
			if err := setFieldValue(fieldVal, valueStr); err != nil {
				return err
			}
		}
	}

	return nil
}

func parseTag(tag string) map[string]string {
	tagParts := strings.Split(tag, ",")
	tagMap := make(map[string]string)

	for _, part := range tagParts {
		if kv := strings.SplitN(part, "=", 2); len(kv) == 2 {
			tagMap[kv[0]] = kv[1]
		} else {
			tagMap[kv[0]] = ""
		}
	}

	return tagMap
}

func deserializeBody(ctx *fasthttp.RequestCtx, fieldVal reflect.Value) error {
	contentType := string(ctx.Request.Header.ContentType())

	switch {
	case strings.HasPrefix(contentType, "application/json"):
		return json.Unmarshal(ctx.Request.Body(), fieldVal.Addr().Interface())
	case strings.HasPrefix(contentType, "multipart/form-data"):
		return parseMultipartForm(ctx, fieldVal.Addr().Interface())
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
		return parseFormURLEncoded(ctx, fieldVal.Addr().Interface())
	case strings.HasPrefix(contentType, "text/plain"):
		fieldVal.SetString(string(ctx.Request.Body()))
	case strings.HasPrefix(contentType, "application/xml"), strings.HasPrefix(contentType, "text/xml"):
		return xml.Unmarshal(ctx.Request.Body(), fieldVal.Addr().Interface())
	case strings.HasPrefix(contentType, "application/octet-stream"):
		return parseOctetStream(ctx, fieldVal.Addr().Interface())
	case strings.HasPrefix(contentType, "text/html"):
		fieldVal.SetString(string(ctx.Request.Body()))
	case strings.HasPrefix(contentType, "application/pdf"):
		return parseBinaryData(ctx, fieldVal.Addr().Interface())
	case strings.HasPrefix(contentType, "application/zip"):
		return parseBinaryData(ctx, fieldVal.Addr().Interface())
	case strings.HasPrefix(contentType, "image/"):
		return parseBinaryData(ctx, fieldVal.Addr().Interface())
	default:
		return InternalServerError{
			Context:     "/api/contexts/DeserializationError",
			Type:        "DeserializationError",
			Title:       "Internal server error",
			Description: "Unsupported content type: " + contentType,
			Violations: []Violation{
				{
					PropertyPath: "body",
					Message:      "Unsupported content type: " + contentType,
				},
			},
		}
	}

	return nil
}

func parseFormURLEncoded(ctx *fasthttp.RequestCtx, dst any) error {
	formData := ctx.PostArgs()
	data := make(map[string][]any)

	formData.VisitAll(func(key, value []byte) {
		data[string(key)] = append(data[string(key)], string(value))
	})

	return mapToStruct(data, dst)
}

func parseMultipartForm(ctx *fasthttp.RequestCtx, dst any) error {
	mr, err := ctx.MultipartForm()
	if err != nil {
		return InternalServerError{
			Context:     "/api/contexts/DeserializationError",
			Type:        "DeserializationError",
			Status:      StatusInternalServerError,
			Title:       "Internal server error",
			Description: "Failed to parse multipart form",
			Violations: []Violation{
				{
					PropertyPath: "body",
					Message:      err.Error(),
				},
			},
		}
	}

	data := make(map[string][]any)

	for key, values := range mr.Value {
		if len(values) > 0 {
			for _, v := range values {
				data[key] = append(data[key], v)
			}
		}
	}

	for key, files := range mr.File {
		if len(files) > 0 {
			for _, fileHeader := range files {
				data[key] = append(data[key], fileHeader)
			}
		}
	}

	return mapToStruct(data, dst)
}

func parseOctetStream(ctx *fasthttp.RequestCtx, dst any) error {
	return parseBinaryData(ctx, dst)
}

func parseBinaryData(ctx *fasthttp.RequestCtx, dst any) error {
	body := ctx.Request.Body()
	fieldVal := reflect.ValueOf(dst).Elem()

	if fieldVal.Kind() == reflect.Slice && fieldVal.Type().Elem().Kind() == reflect.Uint8 {
		fieldVal.SetBytes(body)

		return nil
	}

	return InternalServerError{
		Context:     "/api/contexts/DeserializationError",
		Type:        "DeserializationError",
		Title:       "Internal server error",
		Description: "Unsupported type for binary data",
		Violations: []Violation{
			{
				PropertyPath: "body",
				Message:      "Unsupported type for binary data",
			},
		},
	}
}

func mapToStruct(data map[string][]any, dst any) (err error) {
	dstVal := reflect.ValueOf(dst).Elem()

	if dstVal.Kind() == reflect.Struct {
		for i := 0; i < dstVal.NumField(); i++ {
			field := dstVal.Type().Field(i)
			fieldVal := dstVal.Field(i)
			key := field.Tag.Get("form")

			if key == "" {
				key = field.Name
			}

			if value, exists := data[key]; exists {
				var valueType reflect.Type
				if len(value) > 0 {
					valueType = reflect.TypeOf(value[0])
				}

				if field.Type.Kind() == reflect.Ptr || field.Type.Kind() == reflect.Pointer { //nolint:gocritic
					if field.Type.Elem().Kind() != reflect.Slice && field.Type.Elem().Kind() != reflect.Array {
						err = setFieldValue(fieldVal, value[0])
					}
				} else if field.Type.Kind() != reflect.Slice && field.Type.Kind() != reflect.Array {
					err = setFieldValue(fieldVal, value[0])
				} else {
					err = setFieldValue(fieldVal, value, valueType)
				}

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func setFieldValue(fieldVal reflect.Value, valueStr any, dataValType ...reflect.Type) error {
	switch fieldVal.Kind() {
	case reflect.Ptr:
		if fieldVal.IsNil() {
			fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
		}

		return setFieldValue(fieldVal.Elem(), valueStr)

	case reflect.Struct:
		if str, ok := valueStr.(string); ok {
			return json.Unmarshal([]byte(str), fieldVal.Addr().Interface())
		}

		val := reflect.ValueOf(valueStr)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		fieldVal.Set(val)

	case reflect.String:
		fieldVal.SetString(valueStr.(string))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setIntValue(fieldVal, valueStr.(string))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUintValue(fieldVal, valueStr.(string))

	case reflect.Float32, reflect.Float64:
		return setFloatValue(fieldVal, valueStr.(string))

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(valueStr.(string))
		if err != nil {
			return InternalServerError{
				Context:     "/api/contexts/DeserializationError",
				Type:        "DeserializationError",
				Title:       "Internal server error",
				Description: "Failed to parse bool",
				Violations: []Violation{
					{
						PropertyPath: fieldVal.Type().Name(),
						Message:      err.Error(),
					},
				},
			}
		}

		fieldVal.SetBool(boolValue)

	case reflect.Slice, reflect.Array:
		if fieldVal.Type().Elem().Kind() == reflect.Uint8 {
			fieldVal.SetBytes([]byte(valueStr.(string)))
		} else {
			// create a new slice of the same type as the data
			if len(dataValType) == 0 {
				return InternalServerError{
					Context:     "/api/contexts/DeserializationError",
					Type:        "DeserializationError",
					Title:       "Internal server error",
					Description: "Unsupported slice type " + fieldVal.Type().Elem().Kind().String(),
					Violations: []Violation{
						{
							PropertyPath: fieldVal.Type().Elem().Name(),
							Message:      "Unsupported slice type " + fieldVal.Type().Elem().Kind().String(),
						},
					},
				}
			}

			sliceType := reflect.SliceOf(dataValType[0])
			newSlice := reflect.MakeSlice(sliceType, 0, 0)

			for _, v := range valueStr.([]any) {
				newSlice = reflect.Append(newSlice, reflect.ValueOf(v))
			}

			fieldVal.Set(newSlice)
		}
	case reflect.Interface:
		fieldVal.Set(reflect.ValueOf(valueStr))
	case reflect.Map:
		if fieldVal.Type().Key().Kind() == reflect.String {
			if err := json.Unmarshal([]byte(valueStr.(string)), fieldVal.Addr().Interface()); err != nil {
				return InternalServerError{
					Context:     "/api/contexts/DeserializationError",
					Type:        "DeserializationError",
					Title:       "Internal server error",
					Description: "Failed to unmarshal map",
					Violations: []Violation{
						{
							PropertyPath: fieldVal.Type().Key().Name(),
							Message:      err.Error(),
						},
					},
				}
			}
		} else {
			return InternalServerError{
				Context:     "/api/contexts/DeserializationError",
				Type:        "DeserializationError",
				Title:       "Internal server error",
				Description: "Unsupported map key type " + fieldVal.Type().Key().Kind().String(),
				Violations: []Violation{
					{
						PropertyPath: fieldVal.Type().Key().Name(),
						Message:      "Unsupported map key type " + fieldVal.Type().Key().Kind().String(),
					},
				},
			}
		}

	case reflect.Invalid, reflect.Uintptr, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func,
		reflect.UnsafePointer:
		fallthrough
	default:
		return InternalServerError{
			Context:     "/api/contexts/DeserializationError",
			Type:        "DeserializationError",
			Title:       "Internal server error",
			Description: "Unsupported kind " + fieldVal.Kind().String(),
			Violations: []Violation{
				{
					PropertyPath: fieldVal.Type().Name(),
					Message:      "Unsupported kind " + fieldVal.Kind().String(),
				},
			},
		}
	}

	return nil
}

func setIntValue(fieldVal reflect.Value, valueStr string) error {
	intValue, err := strconv.ParseInt(valueStr, 10, fieldVal.Type().Bits())
	if err != nil {
		return InternalServerError{
			Context:     "/api/contexts/DeserializationError",
			Type:        "DeserializationError",
			Title:       "Internal server error",
			Description: "Failed to parse int",
			Violations: []Violation{
				{
					PropertyPath: fieldVal.Type().Name(),
					Message:      err.Error(),
				},
			},
		}
	}

	fieldVal.SetInt(intValue)

	return nil
}

func setUintValue(fieldVal reflect.Value, valueStr string) error {
	uintValue, err := strconv.ParseUint(valueStr, 10, fieldVal.Type().Bits())
	if err != nil {
		return InternalServerError{
			Context:     "/api/contexts/DeserializationError",
			Type:        "DeserializationError",
			Title:       "Internal server error",
			Description: "Failed to parse uint",
			Violations: []Violation{
				{
					PropertyPath: fieldVal.Type().Name(),
					Message:      err.Error(),
				},
			},
		}
	}

	fieldVal.SetUint(uintValue)

	return nil
}

func setFloatValue(fieldVal reflect.Value, valueStr string) error {
	floatValue, err := strconv.ParseFloat(valueStr, fieldVal.Type().Bits())
	if err != nil {
		return InternalServerError{
			Context:     "/api/contexts/DeserializationError",
			Type:        "DeserializationError",
			Title:       "Internal server error",
			Description: "Failed to parse float",
			Violations: []Violation{
				{
					PropertyPath: fieldVal.Type().Name(),
					Message:      err.Error(),
				},
			},
		}
	}

	fieldVal.SetFloat(floatValue)

	return nil
}

func buildRegex(route string) *regexp.Regexp {
	re := regexp.MustCompile(`:[^/]+`)
	pattern := re.ReplaceAllString(route, `([^/]+)`)
	pattern = "^" + pattern + "$"

	return regexp.MustCompile(pattern)
}

func extractParams(path string, reqPath string) map[string]string {
	var paramNames []string

	for _, segment := range strings.Split(path, "/") {
		if strings.HasPrefix(segment, ":") {
			paramNames = append(paramNames, segment[1:])
		}
	}

	if len(paramNames) == 0 {
		return nil
	}

	re := buildRegex(path)

	matches := re.FindStringSubmatch(reqPath)
	if matches == nil {
		return nil
	}

	params := make(map[string]string)
	for i, name := range paramNames {
		params[name] = matches[i+1]
	}

	return params
}
