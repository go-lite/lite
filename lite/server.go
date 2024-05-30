package lite

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/invopop/yaml"
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

func NewApp() *App {
	return &App{
		App:           fiber.New(),
		OpenApiSpec:   NewOpenApiSpec(),
		OpenAPIConfig: defaultOpenAPIConfig,
	}
}

// AddTags adds tags from the Server (i.e Group)
// Tags from the parent Groups will be respected
func (s *App) AddTags(tags ...string) *App {
	s.tags = append(s.tags, tags...)
	return s
}

// SaveOpenAPISpec saves the OpenAPI spec to a file in YAML format
func (s *App) SaveOpenAPISpec() ([]byte, error) {
	json, err := s.OpenApiSpec.MarshalJSON()
	if err != nil {
		return nil, err
	}

	yamlData, err := writeOpenAPISpec(json)
	if err != nil {
		return nil, err
	}

	return yamlData, nil
}

// writeOpenAPISpec writes the OpenAPI spec to a file in YAML format
func writeOpenAPISpec(d []byte) ([]byte, error) {
	// convert json to yaml
	yamlData, err := yaml.JSONToYAML(d)
	if err != nil {
		return nil, err
	}

	return yamlData, nil
}

func (s *App) AddServer(url, description string) {
	var servers []*openapi3.Server

	servers = append(servers, &openapi3.Server{
		URL:         url,
		Description: description,
	})

	s.OpenApiSpec.Servers = servers
}

// Description sets the description of the OpenAPI spec
func (s *App) Description(description string) *App {
	s.OpenApiSpec.Info.Description = description

	return s
}

// Title sets the title of the OpenAPI spec
func (s *App) Title(title string) *App {
	s.OpenApiSpec.Info.Title = title

	return s
}

// Version sets the version of the OpenAPI spec
func (s *App) Version(version string) *App {
	s.OpenApiSpec.Info.Version = version

	return s
}
