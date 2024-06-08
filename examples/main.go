package main

import (
	"errors"
	"github.com/disco07/lite-fiber"
	"github.com/disco07/lite-fiber/codec"
	"github.com/disco07/lite-fiber/examples/parameters"
	"github.com/disco07/lite-fiber/examples/returns"
	"log"
	"os"
)

// Define example handler
//func getHandler(c *openapi.ContextWithRequest[parameters.GetReq]) (returns.GetResponse, error) {
//	body, err := c.Body()
//	if err != nil {
//		return returns.GetResponse{}, err
//	}
//
//	if body.Params.Value.Name == "test" {
//		return returns.GetResponse{}, errors.New("test is not valid name")
//	}
//
//	return returns.GetResponse{
//		Message: "Hello World!, -" + pathParams.Name + " - " + body.Params.Value.Name,
//	}, nil
//}

func postHandler(c *openapi.ContextWithRequest[parameters.CreateReq]) (returns.CreateResponse, error) {
	body, err := c.Request()
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

func getArrayHandler(_ *openapi.ContextWithRequest[parameters.GetArrayReq]) (returns.GetArrayReturnsResponse, error) {
	res := make([]string, 0)
	res = append(res, "Hello World!")

	return res, nil
}

func main() {
	liteApp := openapi.NewApp()

	//openapi.Get[returns.GetResponse, parameters.GetReq, codec.AsJSON[returns.GetResponse]](liteApp, "/example/:name", getHandler)
	openapi.Post[
		returns.CreateResponse,
		parameters.CreateReq,
		codec.AsJSON[returns.CreateResponse],
	](liteApp, "/example/:id", postHandler).
		OperationID("createExample").
		Description("Create example").
		AddTags("example")

	openapi.Get[returns.GetArrayReturnsResponse, parameters.GetArrayReq, codec.AsJSON[returns.GetArrayReturnsResponse]](liteApp, "/example", getArrayHandler)

	liteApp.AddServer("http://localhost:6000", "example server")

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

	log.Fatal(liteApp.Listen(":6000"))
}
