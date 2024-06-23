package main

import (
	"errors"
	"log"
	"os"

	"github.com/disco07/lite/examples/basic/parameters"
	"github.com/disco07/lite/examples/basic/returns"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/disco07/lite"
)

// Define example handler
func getHandler(c *lite.ContextWithRequest[parameters.GetReq]) (returns.GetResponse, error) {
	request, err := c.Requests()
	if err != nil {
		return returns.GetResponse{}, err
	}

	if request.Params == "test" {
		return returns.GetResponse{}, errors.New("test is not valid name")
	}

	return returns.GetResponse{
		Message: "Hello World!, " + request.Params,
	}, nil
}

func postHandler(c *lite.ContextWithRequest[parameters.CreateReq]) (returns.CreateResponse, error) {
	request, err := c.Requests()
	if err != nil {
		return returns.CreateResponse{}, err
	}

	if request.Body.FirstName == "" {
		return returns.CreateResponse{}, errors.New("first_name are required")
	}

	return returns.CreateResponse{
		ID:        request.ID,
		FirstName: request.Body.FirstName,
		LastName:  request.Body.LastName,
	}, nil
}

func getArrayHandler(_ *lite.ContextWithRequest[parameters.GetArrayReq]) (returns.GetArrayReturnsResponse, error) {
	res := make([]returns.Ret, 0)

	value := "value"
	res = append(res, returns.Ret{
		Message: "Hello World!",
		Embed: returns.Embed{
			Key:        "key",
			ValueEmbed: &value,
		},
	},
		returns.Ret{
			Message: "Hello World 2!",
			Embed: returns.Embed{
				Key: "key2",
			},
		},
	)

	return res, nil
}

func main() {
	app := lite.New()

	app.Use(logger.New())
	app.Use(recover.New())

	lite.Get(app, "/example/:name", getHandler).SetResponseContentType("application/xml")

	lite.Post(app, "/example/:id", postHandler).
		OperationID("createExample").
		Description("Create example").
		AddTags("example")

	lite.Get(app, "/example", getArrayHandler)

	app.AddServer("http://localhost:6001", "example server")

	yamlBytes, err := app.SaveOpenAPISpec()
	if err != nil {
		return
	}

	f, err := os.Create("./examples/basic/api/openapi.yaml")
	if err != nil {
		return
	}

	defer func() {
		closeErr := f.Close()
		if err != nil {
			if closeErr != nil {
				err = closeErr
			}

			log.Fatal(err)
		}
	}()

	_, err = f.Write(yamlBytes)
	if err != nil {
		return
	}

	if err = app.Listen(":6001"); err != nil {
		return
	}
}
