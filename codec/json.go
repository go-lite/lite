package codec

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type AsJSON[T any] struct {
	Value T
}

func (j AsJSON[T]) Decode(c *fiber.Ctx) (T, error) {
	var value T

	err := c.BodyParser(&value)
	if err != nil {
		return value, err
	}

	return value, err
}

func (j AsJSON[T]) Encode(c *fiber.Ctx, value T) error {
	return c.JSON(value)
}

func (j AsJSON[T]) ContentType() string {
	return "application/json"
}

func (j AsJSON[T]) TypeOf() reflect.Type {
	return reflect.TypeOf(j.Value)
}
