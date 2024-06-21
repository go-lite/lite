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
go get github.com/disco07/lite
```

## Usage
Here is a simple example of how to use Lite:

```go
package main

import (
	"github.com/disco07/lite"
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