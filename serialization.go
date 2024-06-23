package lite

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"reflect"

	"github.com/valyala/fasthttp"
)

func serializeResponse(ctx *fasthttp.RequestCtx, src any) error {
	if src == nil {
		return nil
	}

	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	switch srcVal.Kind() {
	case reflect.String:
		ctx.Response.Header.SetContentType("text/plain; charset=utf-8")
		ctx.Response.SetBody([]byte(srcVal.String()))

		return nil
	case reflect.Invalid, reflect.Ptr, reflect.Chan, reflect.Func, reflect.Complex64, reflect.Complex128,
		reflect.Uintptr, reflect.UnsafePointer:
		err := fmt.Errorf("unsupported type: %s", srcVal.Kind())

		ctx.Error(err.Error(), StatusInternalServerError)

		return err
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64,
		reflect.Array, reflect.Interface, reflect.Map, reflect.Slice, reflect.Struct:
		fallthrough
	default:
		return serialize(ctx, srcVal)
	}
}

func serialize(ctx *fasthttp.RequestCtx, srcVal reflect.Value) error {
	contentType := ctx.Response.Header.ContentType()

	switch string(contentType) {
	case "application/json":
		if err := json.NewEncoder(ctx).Encode(srcVal.Interface()); err != nil {
			ctx.Error(err.Error(), StatusInternalServerError)

			return err
		}
	case "application/xml":
		if err := xml.NewEncoder(ctx).Encode(srcVal.Interface()); err != nil {
			ctx.Error(err.Error(), StatusInternalServerError)

			return err
		}
	case "multipart/form-data", "application/x-www-form-urlencoded":
		if form, ok := srcVal.Interface().(map[string]string); ok {
			formData := url.Values{}

			for key, value := range form {
				formData.Set(key, value)
			}

			ctx.SetBody([]byte(formData.Encode()))
		} else {
			err := fmt.Errorf("expected map[string]string for form data serialization")
			ctx.Error(err.Error(), StatusInternalServerError)

			return err
		}
	case "application/octet-stream":
		if data, ok := srcVal.Interface().([]byte); ok {
			ctx.SetBody(data)
		} else {
			err := fmt.Errorf("expected []byte for octet-stream serialization")
			ctx.Error(err.Error(), StatusInternalServerError)

			return err
		}
	case "application/pdf", "application/zip", "image/png", "image/jpeg", "image/gif", "image/webp", "image/svg+xml",
		"image/tiff", "image/vnd.microsoft.icon", "image/vnd.wap.wbmp", "image/x-icon", "image/x-jng", "image/jpg":
		if data, ok := srcVal.Interface().([]byte); ok {
			ctx.SetBody(data)
		} else {
			err := fmt.Errorf("expected []byte for binary file serialization")
			ctx.Error(err.Error(), StatusInternalServerError)

			return err
		}
	default:
		err := fmt.Errorf("unsupported content type: %s", contentType)
		ctx.Error(err.Error(), StatusInternalServerError)

		return err
	}

	return nil
}
