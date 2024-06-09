package openapi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func deserializeParams(ctx *fasthttp.RequestCtx, dst any, params map[string]string) error {
	return deserialize(ctx, reflect.ValueOf(dst).Elem(), params)
}

func deserialize(ctx *fasthttp.RequestCtx, dstVal reflect.Value, params map[string]string) error {
	dstType := dstVal.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		fieldVal := dstVal.Field(i)
		tag := field.Tag.Get("lite")

		if fieldVal.Kind() == reflect.Struct && tag == "" {
			// Recursively handle nested structs
			if err := deserialize(ctx, fieldVal, params); err != nil {
				return err
			}

			continue
		}

		if tag == "" {
			tag = field.Name
		}

		tagMap := parseTag(tag)

		if val, ok := tagMap["req"]; ok && val == "body" {
			continue
		}

		var valueStr string
		if pathKey, ok := tagMap["path"]; ok {
			if params != nil {
				paramsValue, ok := params[pathKey]
				if !ok {
					return fmt.Errorf("missing path parameter: %s", pathKey)
				}

				valueStr = paramsValue
			}
		} else if queryKey, ok := tagMap["query"]; ok {
			if value := ctx.QueryArgs().Peek(queryKey); len(value) > 0 {
				valueStr = string(value)
			}
		} else if headerKey, ok := tagMap["header"]; ok {
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

func deserializeBody(ctx *fasthttp.RequestCtx, dst any) error {
	contentType := string(ctx.Request.Header.ContentType())
	dstVal := reflect.ValueOf(dst).Elem()
	dstType := dstVal.Type()

	for i := 0; i < dstType.NumField(); i++ {
		field := dstType.Field(i)
		fieldVal := dstVal.Field(i)
		tag := field.Tag.Get("lite")
		if tag == "" {
			continue
		}

		tagParts := strings.Split(tag, "=")
		if len(tagParts) != 2 || tagParts[1] != "body" {
			continue
		}

		switch {
		case strings.HasPrefix(contentType, "application/json"):
			return json.Unmarshal(ctx.Request.Body(), fieldVal.Addr().Interface())
		case strings.HasPrefix(contentType, "multipart/form-data"):
			return parseMultipartForm(ctx, fieldVal.Addr().Interface())
		case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
			return parseFormURLEncoded(ctx, fieldVal.Addr().Interface())
		case strings.HasPrefix(contentType, "text/plain"):
			fieldVal.SetString(string(ctx.Request.Body()))
		case strings.HasPrefix(contentType, "application/xml"):
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
			fileHeader := files[0]

			file, err := fileHeader.Open()
			if err != nil {
				return err
			}
			defer file.Close()

			content := new(strings.Builder)

			_, err = io.Copy(content, file)
			if err != nil {
				return err
			}

			data[key] = content.String()
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
		key := field.Tag.Get("json")

		if key == "" {
			key = field.Name
		}

		if value, exists := data[key]; exists {
			valueStr := fmt.Sprintf("%v", value)
			if err := setFieldValue(fieldVal, valueStr); err != nil {
				return err
			}
		}
	}
	return nil
}

func setFieldValue(fieldVal reflect.Value, valueStr string) error {
	switch fieldVal.Kind() {
	case reflect.Ptr:
		if fieldVal.IsNil() {
			fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
		}

		return setFieldValue(fieldVal.Elem(), valueStr)
	case reflect.String:
		fieldVal.SetString(valueStr)
	case reflect.Int:
		intValue, err := strconv.Atoi(valueStr)
		if err != nil {
			return err
		}

		fieldVal.SetInt(int64(intValue))
	case reflect.Int8:
		intValue, err := strconv.ParseInt(valueStr, 10, 8)
		if err != nil {
			return err
		}

		fieldVal.SetInt(intValue)
	case reflect.Int16:
		intValue, err := strconv.ParseInt(valueStr, 10, 16)
		if err != nil {
			return err
		}

		fieldVal.SetInt(intValue)
	case reflect.Int32:
		intValue, err := strconv.ParseInt(valueStr, 10, 32)
		if err != nil {
			return err
		}

		fieldVal.SetInt(intValue)

	case reflect.Int64:
		intValue, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return err
		}

		fieldVal.SetInt(intValue)

	case reflect.Uint:
		uintValue, err := strconv.ParseUint(valueStr, 10, 0)
		if err != nil {
			return err
		}

		fieldVal.SetUint(uintValue)
	case reflect.Uint8:
		uintValue, err := strconv.ParseUint(valueStr, 10, 8)
		if err != nil {
			return err
		}

		fieldVal.SetUint(uintValue)
	case reflect.Uint16:
		uintValue, err := strconv.ParseUint(valueStr, 10, 16)
		if err != nil {
			return err
		}

		fieldVal.SetUint(uintValue)
	case reflect.Uint32:
		uintValue, err := strconv.ParseUint(valueStr, 10, 32)
		if err != nil {
			return err
		}

		fieldVal.SetUint(uintValue)
	case reflect.Uint64:
		uintValue, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			return err
		}

		fieldVal.SetUint(uintValue)
	case reflect.Float32:
		floatValue, err := strconv.ParseFloat(valueStr, 32)
		if err != nil {
			return err
		}

		fieldVal.SetFloat(floatValue)
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return err
		}

		fieldVal.SetFloat(floatValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(valueStr)
		if err != nil {
			return err
		}

		fieldVal.SetBool(boolValue)
	case reflect.Slice:
		if fieldVal.Type().Elem().Kind() == reflect.Uint8 {
			fieldVal.SetBytes([]byte(valueStr))
		} else {
			return fmt.Errorf("unsupported slice type %s", fieldVal.Type().Elem().Kind())
		}
	default:
		return fmt.Errorf("unsupported kind %s", fieldVal.Kind())
	}

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
