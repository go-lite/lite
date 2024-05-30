package lite

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"log"
)

type OpenAPIConfig struct {
	DisableSwagger   bool                               // If true, the server will not serve the swagger ui nor the openapi json spec
	DisableLocalSave bool                               // If true, the server will not save the openapi json spec locally
	SwaggerUrl       string                             // URL to serve the swagger ui
	UIHandler        func(specURL string) fiber.Handler // Handler to serve the openapi ui from spec url
	YamlUrl          string                             // Local path to save the openapi json spec
}

func NewOpenApiSpec() openapi3.T {
	info := &openapi3.Info{
		Title:       "OpenAPI",
		Description: "OpenAPI",
		Version:     "0.0.1",
	}
	spec := openapi3.T{
		OpenAPI: "3.0.3",
		Info:    info,
		Paths:   &openapi3.Paths{},
		Components: &openapi3.Components{
			Schemas:         make(map[string]*openapi3.SchemaRef),
			Parameters:      make(map[string]*openapi3.ParameterRef),
			Headers:         make(map[string]*openapi3.HeaderRef),
			RequestBodies:   make(map[string]*openapi3.RequestBodyRef),
			Responses:       make(map[string]*openapi3.ResponseRef),
			SecuritySchemes: make(map[string]*openapi3.SecuritySchemeRef),
		},
	}
	return spec
}

var defaultOpenAPIConfig = OpenAPIConfig{
	SwaggerUrl: "/swagger",
	YamlUrl:    "/swagger/openapi.yaml",
	UIHandler:  DefaultOpenAPIHandler,
}

type App struct {
	*fiber.App

	OpenApiSpec   openapi3.T
	OpenAPIConfig OpenAPIConfig
	// OpenAPI documentation tags used for logical groupings of operations
	// These tags will be inherited by child Routes/Groups
	tags []string
}

func NewApp(app *fiber.App) *App {
	return &App{
		App:           app,
		OpenApiSpec:   NewOpenApiSpec(),
		OpenAPIConfig: defaultOpenAPIConfig,
	}
}

func (s *App) GetTags() []string {
	return s.tags
}

func (s *App) Tags(tags ...string) *App {
	s.tags = tags
	return s
}

// AddTags adds tags from the Server (i.e Group)
// Tags from the parent Groups will be respected
func (s *App) AddTags(tags ...string) *App {
	s.tags = append(s.tags, tags...)
	return s
}

// SaveOpenAPISpec saves the OpenAPI spec to a file in YAML format
func (s *App) SaveOpenAPISpec() error {
	json, err := s.OpenApiSpec.MarshalJSON()
	if err != nil {
		return err
	}

	log.Println(string(json))

	return nil
}

func (s *App) AddServer(url, description string) {
	var servers []*openapi3.Server

	servers = append(servers, &openapi3.Server{
		URL:         url,
		Description: description,
	})

	s.OpenApiSpec.Servers = servers
}
