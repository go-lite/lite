package codec

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type AsPathParam[T any] struct {
	Value T
}

func (p AsPathParam[T]) Decode(c *fiber.Ctx) (T, error) {
	var value T

	err := c.ParamsParser(&value)
	if err != nil {
		return value, err
	}

	return value, nil
}

func (p AsPathParam[T]) ParamType() string {
	return "path"
}

func (a AsPathParam[T]) TypeOf() reflect.Type {
	return reflect.TypeOf(a.Value)
}
