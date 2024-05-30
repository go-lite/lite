package codec

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type AsHeader[T any] struct {
	Value T
}

func (h AsHeader[T]) Decode(c *fiber.Ctx) (T, error) {
	var value T

	err := c.QueryParser(&value)
	if err != nil {
		return value, err
	}

	return value, nil
}

func (h AsHeader[T]) ParamType() string {
	return "header"
}

func (h AsHeader[T]) TypeOf() reflect.Type {
	return reflect.TypeOf(h.Value)
}
