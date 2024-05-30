package codec

import "github.com/gofiber/fiber/v2"

type AsPDF[T any] struct {
	Value T
}

func (p AsPDF[T]) Decode(c *fiber.Ctx) (T, error) {
	var value T

	err := c.BodyParser(&value)
	if err != nil {
		return value, err
	}

	return value, err
}

func (p AsPDF[T]) Encode(c *fiber.Ctx, value T) error {
	c.Set("Content-Type", "application/pdf")
	return c.SendString("PDF content") // Example, should send actual PDF content
}

func (p AsPDF[T]) ContentType() string {
	return "application/pdf"
}

func (p AsPDF[T]) StructTag() string {
	return "pdf"
}
