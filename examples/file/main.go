package main

import (
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/go-lite/lite"
	"github.com/go-lite/lite/errors"
	"github.com/go-lite/lite/mime"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type ImageResponse = []byte

type ImagePayload struct {
	Body Image `lite:"req=body,multipart/form-data"`
}

type Image struct {
	Info  info                  `form:"info"`
	Image *multipart.FileHeader `form:"image"`
}

type info struct {
	FileName string `form:"filename"`
}

func main() {
	app := lite.New()

	lite.Use(app, logger.New())
	lite.Use(app, recover.New())

	lite.Post(app, "/v1/image/analyse", func(c *lite.ContextWithRequest[ImagePayload]) (ImageResponse, error) {
		req, err := c.Requests()
		if err != nil {
			return ImageResponse{}, errors.NewBadRequestError(err.Error())
		}

		image := req.Body.Image

		if err = c.SaveFile(image, "./examples/file/uploads/"+image.Filename); err != nil {
			return ImageResponse{}, err
		}

		// get the file
		f, err := os.Open("./examples/file/uploads/" + image.Filename)
		if err != nil {
			return ImageResponse{}, err
		}

		// Dummy data for the response
		response, err := io.ReadAll(f)
		if err != nil {
			log.Fatalf("failed reading file: %s", err)
		}

		c.SetContentType(mime.ImagePng)

		return response, nil
	}).SetResponseContentType("image/png")
	lite.Post(app, "/v1/pdf", func(c *lite.ContextWithRequest[[]byte]) (any, error) {
		req, err := c.Requests()
		if err != nil {
			return nil, errors.NewBadRequestError(err.Error())
		}

		log.Println(string(req))

		return nil, nil
	})

	if err := app.Run(); err != nil {
		return
	}
}
