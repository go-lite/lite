package codec

import "github.com/gofiber/fiber/v2"

type AsXML[T any] struct {
	Value T
}

func (x AsXML[T]) Encode(c *fiber.Ctx, value T) error {
	c.Set("Content-Type", "application/xml")

	return c.SendString("application/xml")
}

func (x AsXML[T]) ContentType() string {
	return "application/xml"
}

func (x AsXML[T]) StructTag() string {
	return "xml"
}
