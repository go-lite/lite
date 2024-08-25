package lite

import (
	"encoding/json"
	"encoding/xml"
	"errors"
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
	case reflect.Invalid, reflect.Chan, reflect.Func, reflect.Complex64, reflect.Complex128,
		reflect.Uintptr, reflect.UnsafePointer:
		err := fmt.Errorf("unsupported type: %s", srcVal.Kind())

		ctx.Error(err.Error(), StatusInternalServerError)

		return err
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64,
		reflect.Array, reflect.Interface, reflect.Map, reflect.Slice, reflect.Struct, reflect.Ptr:
		fallthrough
	default:
		return serialize(ctx, srcVal)
	}
}

func serialize(ctx *fasthttp.RequestCtx, srcVal reflect.Value) error {
	contentType := ctx.Response.Header.ContentType()

	switch ContentType(contentType) {
	case ContentTypeJSON:
		if err := json.NewEncoder(ctx).Encode(srcVal.Interface()); err != nil {
			ctx.Error(err.Error(), StatusInternalServerError)

			return err
		}
	case ContentTypeXML:
		if err := xml.NewEncoder(ctx).Encode(srcVal.Interface()); err != nil {
			ctx.Error(err.Error(), StatusInternalServerError)

			return err
		}
	case ContentTypeXFormData, ContentTypeFormData:
		if form, ok := srcVal.Interface().(map[string]string); ok {
			formData := url.Values{}

			for key, value := range form {
				formData.Set(key, value)
			}

			ctx.SetBody([]byte(formData.Encode()))
		} else {
			err := errors.New("expected map[string]string for form data serialization")
			ctx.Error(err.Error(), StatusInternalServerError)

			return err
		}
	case ContentTypeOctetStream, ContentTypePDF, ContentTypeZIP, ContentTypePNG, ContentTypeJPEG, ContentTypeGIF,
		ContentTypeWEBP, ContentTypeSVG, ContentTypeTIFF, ContentTypeICO, ContentTypeJNG, ContentTypeDOC, ContentTypeBMP,
		ContentTypeWOFF, ContentTypeWOFF2, ContentTypeJAR, ContentTypeHQX, ContentTypeXLS, ContentTypeXLSX, ContentTypePPT,
		ContentTypePPTX, ContentTypeDOCX, ContentTypeWMLC, ContentTypeWASM, ContentType7Z, ContentTypeCCO, ContentTypeJARDIFF,
		ContentTypeJNLP, ContentTypeEOT, ContentTypeODG, ContentTypeODP, ContentTypeODS, ContentTypeODT, ContentTypeRAR,
		ContentTypeRPM, ContentTypeSEA, ContentTypeSWF, ContentTypeSIT, ContentTypeTCL, ContentTypeCRT, ContentTypeXPI,
		ContentTypeXHTML, ContentTypeAVIF, ContentTypeWBMP, ContentTypePS, ContentTypeRTF, ContentTypeM3U8, ContentTypeKML,
		ContentTypeKMZ, ContentTypeXSPF, ContentTypeRUN, ContentTypePL, ContentTypePRC:
		if data, ok := srcVal.Interface().([]byte); ok {
			ctx.SetBody(data)
		} else {
			err := errors.New("expected []byte for binary file serialization")
			ctx.Error(err.Error(), StatusInternalServerError)

			return err
		}
	case ContentTypeTXT, ContentTypeHTML, ContentTypeCSS, ContentTypeJS, ContentTypeATOM,
		ContentTypeRSS, ContentTypeMML, ContentTypeJAD, ContentTypeWML, ContentTypeHTC:
		if data, ok := srcVal.Interface().(string); ok {
			ctx.SetBodyString(data)
		} else {
			err := errors.New("expected string for text serialization")
			ctx.Error(err.Error(), StatusInternalServerError)

			return err
		}
	case ContentTypeMIDI, ContentTypeMP3, ContentTypeOGG, ContentTypeM4A, ContentTypeRA,
		ContentType3GP, ContentTypeTS, ContentTypeMP4, ContentTypeMPEG, ContentTypeMOV,
		ContentTypeWEBM, ContentTypeFLV, ContentTypeM4V, ContentTypeMNG, ContentTypeASX,
		ContentTypeWMV, ContentTypeAVI:
		if data, ok := srcVal.Interface().([]byte); ok {
			ctx.SetBody(data)
		} else {
			err := errors.New("expected []byte for binary file serialization")
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
