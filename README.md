<p style="text-align: center;">
  <img src="./logo/lite.png" height="200" alt="Lite Logo" />
</p>

[![Go](https://github.com/go-lite/lite/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/go-lite/lite/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-lite/lite.svg)](https://pkg.go.dev/github.com/go-lite/lite)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-lite/lite)](https://goreportcard.com/report/github.com/go-lite/lite)
[![codecov](https://codecov.io/gh/go-lite/lite/graph/badge.svg?token=5OFXTQKHEE)](https://codecov.io/gh/go-lite/lite)
# Lite: A Typed Wrapper for GoFiber
Lite is a typed wrapper for GoFiber, a web framework for Go. It is designed to be lightweight and easy to use, while still providing a powerful API for building web applications. Lite is built on top of GoFiber, so it inherits all of its features and performance benefits.

## Features
- **Typed Requests**: Define request types to ensure correct data handling.
- **Typed Responses**: Define response types to ensure correct data serialization.
- **Error Handling**: Simplify error management with typed responses.
- **Middleware**: Use middleware to add functionality to your routes.
- **OpenAPI Specification**: Generate OpenAPI specs from your routes.

## Installation
To install Lite, use `go get`:

```bash
go get github.com/go-lite/lite
```

## Usage
### Simple Example
Here is a simple example of how to use Lite:

```go
package main

import (
	"github.com/go-lite/lite"
	"log"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	app := lite.New()

	lite.Get(app, "/", func(c *lite.ContextNoRequest) (Response, error) {
		return Response{Message: "Hello, world!"}, nil
	})

	log.Fatal(app.Listen(":3000"))
}
```
The swagger specs is available at `http://localhost:3000/swagger/index.html` if port `3000` is used.

### Other Examples
```go
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

	app.AddServer("http://localhost:9000", "example server")

	if err := app.Run(); err != nil {
		return
	}
}
```


## Contributing
Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License
Lite is licensed under the MIT License. See [LICENSE](LICENSE) for more information.
[]: # (END)
