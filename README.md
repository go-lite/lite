<p style="text-align: center;">
  <img src="./logo/lite.png" height="200" alt="Fuego Logo" />
</p>

[![Go](https://github.com/go-lite/lite/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/go-lite/lite/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-lite/lite.svg)](https://pkg.go.dev/github.com/go-lite/lite)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-lite/lite)](https://goreportcard.com/report/github.com/go-lite/lite)
[![codecov](https://codecov.io/github/disco07/lite/graph/badge.svg?token=QV9UE6F52R)](https://codecov.io/github/disco07/lite)

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

## Contributing
Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License
Lite is licensed under the MIT License. See [LICENSE](LICENSE) for more information.
[]: # (END)