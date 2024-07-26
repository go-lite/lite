package lite

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-lite/lite/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/invopop/yaml"
)

type TypeOfExtension string

const (
	YAMLExtension TypeOfExtension = "yaml"
	JSONExtension TypeOfExtension = "json"
)

var (
	yamlJSONToYAML = yaml.JSONToYAML
	osMkdirAll     = os.MkdirAll
	osCreate       = os.Create
)

type writeCloser struct {
	io.WriteCloser
}

var wrapperWriteCloser = func(w io.WriteCloser) io.WriteCloser {
	return &writeCloser{w}
}

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

	address string // Address to listen on

	serverURL string
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

func SetOpenAPIPath(path string) Config {
	path = strings.TrimLeft(path, ".")

	return func(s *App) {
		s.openAPIConfig.openapiPath = path
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
func AddTags(tags ...*openapi3.Tag) Config {
	return func(s *App) {
		s.openAPISpec.Tags = append(s.openAPISpec.Tags, tags...)
	}
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

	err := osMkdirAll(jsonFolder, 0o750)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error creating directory %s", jsonFolder))
	}

	f, err := osCreate(path)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error creating file %s", path))
	}

	file := wrapperWriteCloser(f)

	defer file.Close()

	_, err = file.Write(swaggerSpec)
	if err != nil {
		return errors.NewInternalServerError("error writing file")
	}

	return nil
}

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
func AddServer(url, description string) Config {
	return func(s *App) {
		var servers []*openapi3.Server

		s.serverURL = url

		servers = append(servers, &openapi3.Server{
			URL:         url,
			Description: description,
		})

		s.openAPISpec.Servers = servers
	}
}

// SetDescription sets the description of the OpenAPI spec
func SetDescription(description string) Config {
	return func(s *App) {
		s.openAPISpec.Info.Description = description
	}
}

// SetTitle sets the title of the OpenAPI spec
func SetTitle(title string) Config {
	return func(s *App) {
		s.openAPISpec.Info.Title = title
	}
}

// SetVersion sets the version of the OpenAPI spec
func SetVersion(version string) Config {
	return func(s *App) {
		s.openAPISpec.Info.Version = version
	}
}

// SetContact sets the contact of the OpenAPI spec
func SetContact(contact *openapi3.Contact) Config {
	return func(s *App) {
		s.openAPISpec.Info.Contact = contact
	}
}

// SetLicense sets the license of the OpenAPI spec
func SetLicense(license *openapi3.License) Config {
	return func(s *App) {
		s.openAPISpec.Info.License = license
	}
}

// SetTermsOfService sets the terms of service of the OpenAPI spec
func SetTermsOfService(termsOfService string) Config {
	return func(s *App) {
		s.openAPISpec.Info.TermsOfService = termsOfService
	}
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
	if s.serverURL == "" {
		s.serverURL = "http://localhost" + s.address
		s.openAPISpec.Servers = append(s.openAPISpec.Servers, &openapi3.Server{
			URL:         s.serverURL,
			Description: "Local server",
		})
	}

	if s.openAPIConfig.disableSwagger {
		return nil
	}

	if !s.openAPIConfig.disableSwagger {
		s.app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET",
		}))

		// Route to serve the OpenAPI file
		s.app.Get(s.openAPIConfig.openapiPath, s.openAPIPathHandler)

		s.app.Get(s.openAPIConfig.swaggerURL, s.openAPIConfig.uiHandler(s.serverURL+s.openAPIConfig.openapiPath))
	}

	swaggerSpec, err := s.saveOpenAPISpec()
	if err != nil {
		return err
	}

	go func() {
		err := s.saveOpenAPIToFile("."+s.openAPIConfig.openapiPath, swaggerSpec)
		if err != nil {
			slog.ErrorContext(context.Background(), "failed to save openapi spec", slog.Any("error", err))
		}
	}()

	return nil
}

func (s *App) openAPIPathHandler(c *fiber.Ctx) error {
	return c.SendFile("." + s.openAPIConfig.openapiPath)
}

func (s *App) Listen(address string) error {
	s.address = address

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
