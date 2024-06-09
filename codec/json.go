package codec

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type AsJSON[T any] struct {
	Value T
}

func (j AsJSON[T]) Encode(c *fiber.Ctx, value T) error {
	return c.JSON(value)
}

func (j AsJSON[T]) ContentType() string {
	return "application/json"
}

func (j AsJSON[T]) StructTag() string {
	return "json"
}

func (j AsJSON[T]) TypeOf() reflect.Type {
	return reflect.TypeOf(j.Value)
}
