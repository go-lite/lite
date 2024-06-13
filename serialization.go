package lite

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"net/textproto"
	"reflect"
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
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return fmt.Errorf("unsupported type: %s", srcVal.Kind())
	default:
		return serialize(ctx, srcVal)
	}
}

func serialize(ctx *fasthttp.RequestCtx, srcVal reflect.Value) error {
	contentType := ctx.Response.Header.ContentType()

	switch string(contentType) {
	case "application/json":
		if err := json.NewEncoder(ctx).Encode(srcVal.Interface()); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)

			return err
		}
	case "application/xml":
		if err := xml.NewEncoder(ctx).Encode(srcVal.Interface()); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)

			return err
		}
	case "multipart/form-data", "multipart/mixed":
		writer := multipart.NewWriter(ctx)
		for i := 0; i < srcVal.NumField(); i++ {
			fieldVal := srcVal.Field(i)

			partWriter, err := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"application/json"}})
			if err != nil {
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)

				return err
			}

			if err = json.NewEncoder(partWriter).Encode(fieldVal.Interface()); err != nil {
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)

				return err
			}

			if err = writer.Close(); err != nil {
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)

				return err
			}

			ctx.Response.Header.SetContentType(writer.FormDataContentType())
		}
	}

	return nil
}
