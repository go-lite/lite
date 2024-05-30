package main

import (
	"errors"
	"github.com/disco07/lite-fiber/codec"
	"github.com/disco07/lite-fiber/examples/parameters"
	"github.com/disco07/lite-fiber/examples/returns"
	"github.com/disco07/lite-fiber/lite"
	"log"
	"os"
)

// Define example handler
func getHandler(c *lite.ContextWithBody[parameters.GetReq]) (returns.GetResponse, error) {
	body, err := c.Body()
	if err != nil {
		return returns.GetResponse{}, err
	}

	if body.Params.Value.Name == "test" {
		return returns.GetResponse{}, errors.New("test is not valid name")
	}

	return returns.GetResponse{
		Message: "Hello World!, " + body.Params.Value.Name,
	}, nil
}

func postHandler(c *lite.ContextWithBody[parameters.CreateReq]) (returns.CreateResponse, error) {
	body, err := c.Body()
	if err != nil {
		return returns.CreateResponse{}, err
	}

	if body.Body.Value.FirstName == "" {
		return returns.CreateResponse{}, errors.New("first_name are required")
	}

	return returns.CreateResponse{
		ID:        body.Params.Value.ID,
		FirstName: body.Body.Value.FirstName,
		LastName:  *body.Body.Value.LastName,
	}, nil
}

func main() {
	liteApp := lite.NewApp()

	lite.Get[returns.GetResponse, parameters.GetReq, codec.AsJSON[returns.GetResponse]](liteApp, "/example/:name", getHandler)
	lite.Post[
		returns.CreateResponse,
		parameters.CreateReq,
		codec.AsJSON[returns.CreateResponse],
	](liteApp, "/example/:id", postHandler).
		OperationID("createExample").
		Description("Create example").
		AddTags("example")

	liteApp.AddServer("http://localhost:3000", "example server")

	yamlBytes, err := liteApp.SaveOpenAPISpec()
	if err != nil {
		log.Fatal(err)
	}

	// Ouvrir le fichier pour écriture
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

	// Écrire le nouveau contenu dans le fichier
	_, err = f.Write(yamlBytes)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(liteApp.Listen(":3000"))
}
