package codec

import "github.com/gofiber/fiber/v2"

type Encoder[T any] interface {
	ContentType
	Encode(c *fiber.Ctx, value T) error
}
