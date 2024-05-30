package codec

import "github.com/gofiber/fiber/v2"

type AsMultiPart[T any] struct {
	Value T
}

func (m AsMultiPart[T]) Decode(c *fiber.Ctx) (T, error) {
	var value T

	err := c.BodyParser(&value)
	if err != nil {
		return value, err
	}

	return value, err
}

func (m AsMultiPart[T]) Encode(c *fiber.Ctx, value T) error {
	c.Set("Content-Type", "multipart/form-data")

	return c.SendString("Multipart content") // Example, should send actual multipart content
}

func (m AsMultiPart[T]) ContentType() string {
	return "multipart/form-data"
}
