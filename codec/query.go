package codec

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type AsQueryParam[T any] struct {
	Value T
}

func (q AsQueryParam[T]) Decode(c *fiber.Ctx) (T, error) {
	var value T

	err := c.QueryParser(&value)
	if err != nil {
		return value, err
	}

	return value, nil
}

func (q AsQueryParam[T]) ParamType() string {
	return "query"
}

func (a AsQueryParam[T]) TypeOf() reflect.Type {
	return reflect.TypeOf(a.Value)
}
