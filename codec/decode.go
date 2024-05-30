package codec

import "github.com/gofiber/fiber/v2"

type Decoder[T any] interface {
	Decode(c *fiber.Ctx) (T, error)
}
