package main

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"

	"github.com/go-lite/lite"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type GetReq struct {
	Login  string  `lite:"header=Basic,isauth,scheme=basic,name=Basic"`
	Name   string  `lite:"header=name"`
	Value  *string `lite:"header=value"`
	Params string  `lite:"path=params"`
}

type GetArrayReq struct{}

type CreateBody struct {
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type CreateReq struct {
	Authorization *string    `lite:"header=Authorization,isauth,scheme=bearer"`
	ID            uint64     `lite:"path=id"`
	Body          CreateBody `lite:"req=body"`
}

type PutReq struct {
	ID   uint64  `lite:"path=id"`
	Body PutBody `lite:"req=body"`
}

type PutBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GetArrayReturnsResponse = lite.List[Ret]

type Ret struct {
	Message  string                 `json:"message"`
	Embed    Embed                  `json:"embed"`
	Map      map[string]string      `json:"map"`
	OtherMap map[string]OtherEmbed2 `json:"other_map"`
}

type Embed struct {
	Key        string     `json:"key"`
	ValueEmbed ValueEmbed `json:"value"`
	Others     []*string  `json:"others"`
	OtherEmbed OtherEmbed `json:"other_embed"`
}

type ValueEmbed = *string

type OtherEmbed struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OtherEmbed2 struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetResponse struct {
	Message string `json:"message" xml:"message"`
}

type CreateResponse struct {
	ID        uint64  `json:"id"`
	FirstName string  `json:"fist_name"`
	LastName  *string `json:"last_name"`
}

type PutResponse struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Define example handler
func getHandler(c *lite.ContextWithRequest[GetReq]) (GetResponse, error) {
	request, err := c.Requests()
	if err != nil {
		return GetResponse{}, err
	}

	if request.Params == "test" {
		return GetResponse{}, errors.New("test is not valid name")
	}

	return GetResponse{
		Message: "Hello World!, " + request.Params,
	}, nil
}

func postHandler(c *lite.ContextWithRequest[CreateReq]) (CreateResponse, error) {
	request, err := c.Requests()
	if err != nil {
		return CreateResponse{}, err
	}

	if request.Body.FirstName == "" {
		return CreateResponse{}, errors.New("first_name are required")
	}

	return CreateResponse{
		ID:        request.ID,
		FirstName: request.Body.FirstName,
		LastName:  request.Body.LastName,
	}, nil
}

func getArrayHandler(_ *lite.ContextWithRequest[GetArrayReq]) (GetArrayReturnsResponse, error) {
	res := make([]Ret, 0)

	value := "value"
	res = append(res, Ret{
		Message: "Hello World!",
		Embed: Embed{
			Key:        "key",
			ValueEmbed: &value,
		},
	},
		Ret{
			Message: "Hello World 2!",
			Embed: Embed{
				Key: "key2",
			},
		},
	)

	return lite.NewList(res), nil
}

func putHandler(c *lite.ContextWithRequest[PutReq]) (PutResponse, error) {
	request, err := c.Requests()
	if err != nil {
		return PutResponse{}, err
	}

	return PutResponse{
		ID:        request.ID,
		FirstName: request.Body.FirstName,
		LastName:  request.Body.LastName,
	}, nil
}

func main() {
	app := lite.New()

	lite.Use(app, logger.New())
	lite.Use(app, recover.New())

	lite.Get(app, "/example/:params", getHandler).SetResponseContentType("application/xml")

	lite.Post(app, "/example/:id", postHandler).
		OperationID("createExample").
		Description("Create example").
		AddTags("example")

	lite.Get(app, "/example", getArrayHandler)

	lite.Put(app, "/example/:id", putHandler)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Route pour servir le fichier OpenAPI
	app.Get("/api/openapi.yaml", func(c *fiber.Ctx) error {
		cwd, err := os.Getwd()
		if err != nil {
			log.Println("Error getting current working directory:", err)
			return c.Status(500).SendString("Internal Server Error")
		}
		log.Println("Current working directory:", cwd)
		return c.SendFile("./examples/basic/api/openapi.yaml")
	})

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
