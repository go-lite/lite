package lite

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-lite/lite/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/invopop/yaml"
	"github.com/valyala/fasthttp"
)

type TypeOfExtension string

const (
	YAMLExtension TypeOfExtension = "yaml"
	JSONExtension TypeOfExtension = "json"
)

type OpenAPIConfig struct {
	DisableSwagger   bool                               // If true, the server will not serve the swagger ui nor the openapi json spec
	DisableLocalSave bool                               // If true, the server will not save the openapi json spec locally
	SwaggerURL       string                             // URL to serve the swagger ui
	UIHandler        func(specURL string) fiber.Handler // Handler to serve the openapi ui from spec url
	YamlURL          string                             // Local path to save the openapi json spec
	TypeOfExtension  TypeOfExtension                    // Type of extension to use for the openapi spec
}

func NewOpenAPISpec() openapi3.T {
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
	SwaggerURL:      "/swagger/*",
	YamlURL:         "/api/openapi.yaml",
	UIHandler:       DefaultOpenAPIHandler,
	TypeOfExtension: YAMLExtension,
}

type App struct {
	*fiber.App

	openAPISpec   openapi3.T
	openAPIConfig OpenAPIConfig

	Serializer func(ctx *fasthttp.RequestCtx, response any) error

	tag string

	basePath string

	// OpenAPI documentation tags used for logical groupings of operations
	// These tags will be inherited by child Routes/Groups
	tags []string
}

func New() *App {
	return &App{
		App:           fiber.New(),
		openAPISpec:   NewOpenAPISpec(),
		openAPIConfig: defaultOpenAPIConfig,
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
	json, err := s.openAPISpec.MarshalJSON()
	if err != nil {
		return nil, err
	}

	yamlData, err := writeOpenAPISpec(json)
	if err != nil {
		return nil, err
	}

	return yamlData, nil
}

var yamlJSONToYAML = yaml.JSONToYAML

// writeOpenAPISpec writes the OpenAPI spec to a file in YAML format
func writeOpenAPISpec(d []byte) ([]byte, error) {
	// convert json to yaml
	yamlData, err := yamlJSONToYAML(d)
	if err != nil {
		return nil, err
	}

	return yamlData, nil
}

// AddServer adds a server to the OpenAPI spec
func (s *App) AddServer(url, description string) {
	var servers []*openapi3.Server

	servers = append(servers, &openapi3.Server{
		URL:         url,
		Description: description,
	})

	s.openAPISpec.Servers = servers
}

// Description sets the description of the OpenAPI spec
func (s *App) Description(description string) *App {
	s.openAPISpec.Info.Description = description

	return s
}

// Title sets the title of the OpenAPI spec
func (s *App) Title(title string) *App {
	s.openAPISpec.Info.Title = title

	return s
}

// Version sets the version of the OpenAPI spec
func (s *App) Version(version string) *App {
	s.openAPISpec.Info.Version = version

	return s
}

func (s *App) createDefaultErrorResponses() (map[int]*openapi3.Response, error) {
	responses := make(map[int]*openapi3.Response)

	for _, errResponse := range errors.DefaultErrorResponses {
		responseSchema, ok := s.openAPISpec.Components.Schemas["httpGenericError"]
		if !ok {
			var err error

			responseSchema, err = generatorNewSchemaRefForValue(new(errors.HTTPError), s.openAPISpec.Components.Schemas)
			if err != nil {
				return nil, err
			}

			s.openAPISpec.Components.Schemas["httpGenericError"] = responseSchema
		}

		response := openapi3.NewResponse().WithDescription(errResponse.Description())

		var consume []string
		consume = append(consume, errors.DefaultErrorContentTypeResponses...)

		if responseSchema != nil {
			content := openapi3.NewContentWithSchemaRef(
				openapi3.NewSchemaRef(fmt.Sprintf(
					"#/components/schemas/%s",
					"httpGenericError",
				), &openapi3.Schema{}),
				consume,
			)
			response.WithContent(content)
		}

		responses[errResponse.StatusCode()] = response
	}

	return responses, nil
}

func (s *App) Listen(address string) error {
	return s.App.Listen(address)
}

func (s *App) Shutdown() error {
	return s.App.Shutdown()
}
