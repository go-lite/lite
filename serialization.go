package lite

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"net/textproto"
	"reflect"
	"strings"
)

func serializeResponse(ctx *fasthttp.RequestCtx, src any) error {
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	if srcVal.Kind() == reflect.Map || srcVal.Kind() == reflect.Slice {
		tagType := "json"
		return serializeMapOrSlice(ctx, srcVal, tagType)
	}

	srcType := srcVal.Type()
	tagType := ""
	tagMap := make(map[string]string)

	for i := 0; i < srcType.NumField(); i++ {
		field := srcType.Field(i)
		tag := field.Tag.Get("lite")
		if tag == "" {
			continue
		}

		tagParts := strings.Split(tag, ":")
		if len(tagParts) != 2 {
			return fmt.Errorf("invalid tag format: %s", tag)
		}

		if tagType == "" {
			tagType = tagParts[1]
		} else if tagType != tagParts[1] {
			tagType = "multipart"
		}

		tagMap[field.Name] = tagParts[1]
	}

	if tagType == "multipart" {
		return serializeMultipartResponse(ctx, srcVal, tagMap)
	}

	return serializeSinglePartResponse(ctx, srcVal, tagType)
}

func serializeMapOrSlice(ctx *fasthttp.RequestCtx, srcVal reflect.Value, tagType string) error {
	var (
		contentType  string
		responseBody []byte
		err          error
	)

	switch tagType {
	case "json":
		contentType = "application/json"
		responseBody, err = json.Marshal(srcVal.Interface())
	case "xml":
		contentType = "application/xml"
		responseBody, err = xml.Marshal(srcVal.Interface())
	default:
		return fmt.Errorf("unsupported content type: %s", tagType)
	}

	if err != nil {
		return err
	}

	ctx.Response.Header.SetContentType(contentType)
	ctx.Response.SetBody(responseBody)
	return nil
}

func serializeSinglePartResponse(ctx *fasthttp.RequestCtx, srcVal reflect.Value, tagType string) error {
	var (
		contentType  string
		responseBody []byte
		err          error
	)

	switch tagType {
	case "json":
		contentType = "application/json"
		responseBody, err = json.Marshal(srcVal.Interface())
	case "xml":
		contentType = "application/xml"
		responseBody, err = xml.Marshal(srcVal.Interface())
	default:
		return fmt.Errorf("unsupported content type: %s", tagType)
	}

	if err != nil {
		return err
	}

	ctx.Response.Header.SetContentType(contentType)
	ctx.Response.SetBody(responseBody)
	return nil
}

func serializeMultipartResponse(ctx *fasthttp.RequestCtx, srcVal reflect.Value, tagMap map[string]string) error {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	for i := 0; i < srcVal.NumField(); i++ {
		field := srcVal.Type().Field(i)
		fieldVal := srcVal.Field(i)

		tagType, ok := tagMap[field.Name]
		if !ok {
			continue
		}

		var partWriter io.Writer
		var err error

		switch tagType {
		case "json":
			partWriter, err = writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"application/json"}})
			if err != nil {
				return err
			}
			if err := json.NewEncoder(partWriter).Encode(fieldVal.Interface()); err != nil {
				return err
			}
		case "xml":
			partWriter, err = writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"application/xml"}})
			if err != nil {
				return err
			}
			if err := xml.NewEncoder(partWriter).Encode(fieldVal.Interface()); err != nil {
				return err
			}
		case "jpg", "jpeg", "png", "gif":
			partWriter, err = writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"image/" + tagType}})
			if err != nil {
				return err
			}
			if _, err = partWriter.Write(fieldVal.Bytes()); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported part content type: %s", tagType)
		}
	}

	if err := writer.Close(); err != nil {
		return err
	}

	ctx.Response.Header.SetContentType(writer.FormDataContentType())
	ctx.Response.SetBody(b.Bytes())
	return nil
}
