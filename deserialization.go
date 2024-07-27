package lite

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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
			return fmt.Errorf("missing tag for field %s", field.Name)
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
					return fmt.Errorf("missing params parameter: %s", paramsKey)
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
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	return nil
}

func parseFormURLEncoded(ctx *fasthttp.RequestCtx, dst any) error {
	formData := ctx.PostArgs()
	data := make(map[string]any)

	formData.VisitAll(func(key, value []byte) {
		data[string(key)] = string(value)
	})

	return mapToStruct(data, dst)
}

func parseMultipartForm(ctx *fasthttp.RequestCtx, dst any) error {
	mr, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	data := make(map[string]any)

	for key, values := range mr.Value {
		if len(values) > 0 {
			data[key] = values[0]
		}
	}

	for key, files := range mr.File {
		if len(files) > 0 {
			if fileHeader := files[0]; fileHeader != nil {
				data[key] = fileHeader
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

	return fmt.Errorf("unsupported type for binary data")
}

func mapToStruct(data map[string]any, dst any) error {
	dstVal := reflect.ValueOf(dst).Elem()
	for i := 0; i < dstVal.NumField(); i++ {
		field := dstVal.Type().Field(i)
		fieldVal := dstVal.Field(i)
		key := field.Tag.Get("form")

		if key == "" {
			key = field.Name
		}

		if value, exists := data[key]; exists {
			if err := setFieldValue(fieldVal, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func setFieldValue(fieldVal reflect.Value, valueStr any) error {
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
			return err
		}

		fieldVal.SetBool(boolValue)

	case reflect.Slice, reflect.Array:
		if fieldVal.Type().Elem().Kind() == reflect.Uint8 {
			fieldVal.SetBytes([]byte(valueStr.(string)))
		} else {
			return fmt.Errorf("unsupported slice type %s", fieldVal.Type().Elem().Kind())
		}
	case reflect.Interface:
		fieldVal.Set(reflect.ValueOf(valueStr))
	case reflect.Map:
		if fieldVal.Type().Key().Kind() == reflect.String {
			if err := json.Unmarshal([]byte(valueStr.(string)), fieldVal.Addr().Interface()); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported map key type %s", fieldVal.Type().Key().Kind())
		}

	case reflect.Invalid, reflect.Uintptr, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func,
		reflect.UnsafePointer:
		fallthrough
	default:
		return fmt.Errorf("unsupported kind %s", fieldVal.Kind())
	}

	return nil
}

func setIntValue(fieldVal reflect.Value, valueStr string) error {
	intValue, err := strconv.ParseInt(valueStr, 10, fieldVal.Type().Bits())
	if err != nil {
		return err
	}

	fieldVal.SetInt(intValue)

	return nil
}

func setUintValue(fieldVal reflect.Value, valueStr string) error {
	uintValue, err := strconv.ParseUint(valueStr, 10, fieldVal.Type().Bits())
	if err != nil {
		return err
	}

	fieldVal.SetUint(uintValue)

	return nil
}

func setFloatValue(fieldVal reflect.Value, valueStr string) error {
	floatValue, err := strconv.ParseFloat(valueStr, fieldVal.Type().Bits())
	if err != nil {
		return err
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
