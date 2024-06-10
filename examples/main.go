package main

import (
	"errors"
	"log"
	"os"

	"github.com/disco07/lite"
	"github.com/disco07/lite/examples/parameters"
	"github.com/disco07/lite/examples/returns"
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
	body, err := c.Requests()
	if err != nil {
		return returns.CreateResponse{}, err
	}

	if body.Body.FirstName == "" {
		return returns.CreateResponse{}, errors.New("first_name are required")
	}

	return returns.CreateResponse{
		ID:        body.Params.ID,
		FirstName: body.Body.FirstName,
		LastName:  body.Body.LastName,
	}, nil
}

func getArrayHandler(_ *lite.ContextWithRequest[parameters.GetArrayReq]) (returns.GetArrayReturnsResponse, error) {
	res := make([]string, 0)
	res = append(res, "Hello World!")

	return res, nil
}

func main() {
	liteApp := lite.NewApp()

	lite.Get(liteApp, "/example/:name", getHandler)

	lite.Post(liteApp, "/example/:id", postHandler).
		OperationID("createExample").
		Description("Create example").
		AddTags("example")

	lite.Get(liteApp, "/example", getArrayHandler)

	liteApp.AddServer("http://localhost:6000", "example server")

	yamlBytes, err := liteApp.SaveOpenAPISpec()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("./examples/api/openapi.yaml")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	_, err = f.Write(yamlBytes)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(liteApp.Listen(":6000"))
}
