package lite

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type Context[B any] interface {
	Body() (B, error)
}

type ContextWithBody[B any] struct {
	Ctx *fiber.Ctx
}

func (c *ContextWithBody[B]) Body() (B, error) {
	var body B
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
