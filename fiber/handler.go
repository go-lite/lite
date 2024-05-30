package fiber

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type Handler[Req, Res any] func(*Ctx[Req, Res]) error

type ctx[B any] interface{}

func Get2[T, B any, Contexted ctx[B]](app fiber.Router, path string, controller func(Contexted) (T, error), tags []string) {
}

type Ctx[Req, Res any] struct {
	*fiber.Ctx
	Req Req
	Res Res
}

func decode[Req any](c *fiber.Ctx) (Req, error) {
	var req Req
	reqValue := reflect.ValueOf(&req).Elem()
	reqType := reqValue.Type()

	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		fieldValue := reqValue.Field(i)

		if field.Type == reflect.TypeOf(AsHeader[string]{}) {
			headerName := field.Tag.Get("header")
			if headerName == "" {
				headerName = field.Name
			}
			fieldValue.Set(reflect.ValueOf(AsHeader[string]{Value: c.Get(headerName)}))
		} else if field.Type == reflect.TypeOf(AsPathParam[uint64]{}) {
			paramName := field.Tag.Get("param")
			if paramName == "" {
				paramName = field.Name
			}
			value, _ := c.ParamsInt(paramName, 0)
			fieldValue.Set(reflect.ValueOf(AsPathParam[uint64]{Value: uint64(value)}))
		} else if field.Type == reflect.TypeOf(AsQueryParam[string]{}) {
			queryName := field.Tag.Get("query")
			if queryName == "" {
				queryName = field.Name
			}
			fieldValue.Set(reflect.ValueOf(AsQueryParam[string]{Value: c.Query(queryName)}))
		} else if field.Type == reflect.TypeOf(AsJSON[any]{}) {
			if err := c.BodyParser(fieldValue.Addr().Interface()); err != nil {
				return req, err
			}
		}
	}

	return req, nil
}

func encode[Res any](c *fiber.Ctx, res Res) error {
	resValue := reflect.ValueOf(&res).Elem()
	resType := resValue.Type()

	for i := 0; i < resType.NumField(); i++ {
		field := resType.Field(i)
		fieldValue := resValue.Field(i)

		if field.Type == reflect.TypeOf(AsJSON[any]{}) {
			return c.JSON(fieldValue.Interface())
		} else if field.Type == reflect.TypeOf(AsPDF[any]{}) {
			// Handle PDF encoding
			// This is a placeholder, actual implementation would depend on your requirements
			c.Set(fiber.HeaderContentType, "application/pdf")
			return c.Send(fieldValue.Bytes())
		} else if field.Type == reflect.TypeOf(AsMultiPart[any]{}) {
			// Handle MultiPart encoding
			// This is a placeholder, actual implementation would depend on your requirements
		} else if field.Type == reflect.TypeOf(AsTextPlain[string]{}) {
			c.Set(fiber.HeaderContentType, "text/plain")
			return c.SendString(fieldValue.String())
		}
	}

	return nil
}

func customHandler[Req, Res any](handler Handler[Req, Res]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Decode the request
		req, err := decode[Req](c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Create the context
		ctx := &Ctx[Req, Res]{
			Ctx: c,
			Req: req,
		}

		// Call the handler
		if err := handler(ctx); err != nil {
			return err
		}

		// Encode the response
		return encode(c, ctx.Res)
	}
}
