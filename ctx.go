package lite

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type Context[Request any] interface {
	Body() (Request, error)
	Context() context.Context
	Request() (Request, error)
}

type ContextWithRequest[Request any] struct {
	Ctx  *fiber.Ctx
	path string
}

func (c *ContextWithRequest[Request]) Body() (Request, error) {
	var body Request
	v := reflect.ValueOf(&body).Elem()
	t := reflect.TypeOf(&body).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		fieldValue := field.Interface()

		method := reflect.ValueOf(fieldValue).MethodByName("Decode")
		if method.IsValid() {
			results := method.Call([]reflect.Value{reflect.ValueOf(c.Ctx)})

			if len(results) == 2 && !results[1].IsNil() {
				err := results[1].Interface().(error)
				return body, err
			}

			if field.FieldByName("Value").IsValid() && field.FieldByName("Value").CanSet() {
				field.FieldByName("Value").Set(reflect.ValueOf(results[0].Interface()))
			}
		}
	}
	return body, nil
}

func (c *ContextWithRequest[Request]) Context() context.Context {
	return c.Ctx.UserContext()
}

func (c *ContextWithRequest[Request]) Request() (Request, error) {
	var req Request

	typeOfReq := reflect.TypeOf(&req).Elem()

	reqContext := c.Ctx.Context()

	params := extractParams(c.path, string(reqContext.Path()))

	switch typeOfReq.Kind() {
	case reflect.Struct:
		err := deserializeParams(reqContext, &req, params)
		if err != nil {
			return req, err
		}

		err = deserializeBody(reqContext, &req)
		if err != nil {
			return req, err
		}
	default:
		return req, errors.New("unsupported type")
	}

	return req, nil
}
