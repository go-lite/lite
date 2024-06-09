package codec

import "github.com/gofiber/fiber/v2"

type AsMultiPart[T any] struct {
	Value T
}

func (m AsMultiPart[T]) Encode(c *fiber.Ctx, value T) error {
	c.Set("Content-Type", "multipart/form-data")

	return c.SendString("Multipart content") // Example, should send actual multipart content
}

func (m AsMultiPart[T]) ContentType() string {
	return "multipart/form-data"
}

func (m AsMultiPart[T]) StructTag() string {
	return "form"
}
