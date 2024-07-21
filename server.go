package lite

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-lite/lite/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/invopop/yaml"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TypeOfExtension string

const (
	YAMLExtension TypeOfExtension = "yaml"
	JSONExtension TypeOfExtension = "json"
)

type config struct {
	disableSwagger   bool                               // If true, the server will not serve the swagger ui nor the openapi json spec
	disableLocalSave bool                               // If true, the server will not save the openapi json spec locally
	swaggerURL       string                             // URL to serve the swagger ui
	uiHandler        func(specURL string) fiber.Handler // Handler to serve the openapi ui from spec url
	openapiPath      string                             // Local path to save the openapi json spec
	typeOfExtension  TypeOfExtension                    // Type of extension to use for the openapi spec
}

func newOpenAPISpec() openapi3.T {
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

var defaultOpenAPIConfig = config{
	swaggerURL:      "/swagger/*",
	openapiPath:     "/api/openapi.yaml",
	uiHandler:       defaultOpenAPIHandler,
	typeOfExtension: YAMLExtension,
}

type App struct {
	app *fiber.App

	openAPISpec   openapi3.T
	openAPIConfig config

	tag string

	basePath string

	// OpenAPI documentation tags used for logical groupings of operations
	// These tags will be inherited by child Routes/Groups
	tags []string

	address string // Address to listen on
}

func New(config ...Config) *App {
	app := &App{
		app:           fiber.New(),
		openAPISpec:   newOpenAPISpec(),
		openAPIConfig: defaultOpenAPIConfig,
		address:       ":9000",
	}

	for _, c := range config {
		c(app)
	}

	return app
}

type Config func(s *App)

func SetDisableSwagger(disable bool) Config {
	return func(s *App) {
		s.openAPIConfig.disableSwagger = disable
	}
}

func SetDisableLocalSave(disable bool) Config {
	return func(s *App) {
		s.openAPIConfig.disableLocalSave = disable
	}
}

func SetSwaggerURL(url string) Config {
	return func(s *App) {
		s.openAPIConfig.swaggerURL = url
	}
}

func SetUIHandler(handler func(specURL string) fiber.Handler) Config {
	return func(s *App) {
		s.openAPIConfig.uiHandler = handler
	}
}

func SetOpenAPIURL(url string) Config {
	return func(s *App) {
		s.openAPIConfig.openapiPath = url
	}
}

func SetTypeOfExtension(extension TypeOfExtension) Config {
	return func(s *App) {
		s.openAPIConfig.typeOfExtension = extension
	}
}

func SetAddress(address string) Config {
	// check if address is valid
	if !strings.HasPrefix(address, ":") {
		panic("address must start with :")
	}

	port := strings.TrimPrefix(address, ":")

	_, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}

	return func(s *App) {
		s.address = address
	}
}

// AddTags adds tags from the Server (i.e Group)
// Tags from the parent Groups will be respected
func (s *App) AddTags(tags ...string) *App {
	s.tags = append(s.tags, tags...)
	return s
}

// saveOpenAPISpec saves the OpenAPI spec to a file in YAML format
func (s *App) saveOpenAPISpec() ([]byte, error) {
	json, err := s.openAPISpec.MarshalJSON()
	if err != nil {
		return nil, err
	}

	if s.openAPIConfig.typeOfExtension == JSONExtension {
		return json, nil
	}

	yamlData, err := writeOpenAPISpec(json)
	if err != nil {
		return nil, err
	}

	return yamlData, nil
}

func (s *App) saveOpenAPIToFile(path string, swaggerSpec []byte) error {
	jsonFolder := filepath.Dir(path)

	err := os.MkdirAll(jsonFolder, 0o750)
	if err != nil {
		return errors.NewInternalServerError("error creating directory")
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.NewInternalServerError("error creating file")
	}
	defer f.Close()

	_, err = f.Write(swaggerSpec)
	if err != nil {
		return errors.NewInternalServerError("error writing file")
	}

	return nil
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

func (s *App) setup() error {
	if s.openAPIConfig.disableSwagger {
		return nil
	}

	swaggerSpec, err := s.saveOpenAPISpec()
	if err != nil {
		return err
	}

	go func() {
		err := s.saveOpenAPIToFile(s.openAPIConfig.openapiPath, swaggerSpec)
		if err != nil {
			slog.ErrorContext(context.Background(), "failed to save openapi spec", slog.Any("error", err))
		}
	}()

	return nil
}

func (s *App) Listen(address string) error {
	err := s.setup()
	if err != nil {
		return err
	}

	return s.app.Listen(address)
}

func (s *App) Run() error {
	err := s.setup()
	if err != nil {
		return err
	}

	return s.app.Listen(s.address)
}

func (s *App) Shutdown() error {
	return s.app.Shutdown()
}
