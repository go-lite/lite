package lite

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/invopop/yaml"
	"github.com/stretchr/testify/assert"
)

func TestNewOpenAPISpec(t *testing.T) {
	spec := newOpenAPISpec()

	assert.Equal(t, "OpenAPI", spec.Info.Title)
	assert.Equal(t, "OpenAPI", spec.Info.Description)
	assert.Equal(t, "0.0.1", spec.Info.Version)
	assert.Equal(t, "3.0.3", spec.OpenAPI)
	assert.NotNil(t, spec.Paths)
	assert.NotNil(t, spec.Components.Schemas)
	assert.NotNil(t, spec.Components.Parameters)
	assert.NotNil(t, spec.Components.Headers)
	assert.NotNil(t, spec.Components.RequestBodies)
	assert.NotNil(t, spec.Components.Responses)
	assert.NotNil(t, spec.Components.SecuritySchemes)
}

func TestNewApp(t *testing.T) {
	app := New(
		SetDisableSwagger(true),
		SetLogger(slog.Default()),
		SetSwaggerURL("/swagger/*"),
		SetOpenAPIPath("/api/openapi.json"),
		SetUIHandler(defaultOpenAPIHandler),
		SetDisableLocalSave(false),
		SetTypeOfExtension(JSONExtension),
		SetValidator(nil),
		SetAddress(":8080"),
	)

	assert.NotNil(t, app.app)
	assert.Equal(t, "OpenAPI", app.openAPISpec.Info.Title)
	assert.Equal(t, defaultOpenAPIConfig.swaggerURL, app.openAPIConfig.swaggerURL)
	assert.Equal(t, "/api/openapi.json", app.openAPIConfig.openapiPath)
}

func TestApp_AddTags(t *testing.T) {
	app := New(AddTags(&openapi3.Tag{Name: "tag1"}, &openapi3.Tag{Name: "tag2"}))

	assert.Contains(t, app.openAPISpec.Tags, &openapi3.Tag{Name: "tag1"})
	assert.Contains(t, app.openAPISpec.Tags, &openapi3.Tag{Name: "tag2"})
}

func TestApp_SaveOpenAPISpec(t *testing.T) {
	app := New()
	yamlData, err := app.saveOpenAPISpec()

	assert.Nil(t, err)
	assert.NotNil(t, yamlData)
}

func TestApp_AddServer(t *testing.T) {
	app := New(
		AddServer("http://localhost", "Local server"),
	)

	assert.Len(t, app.openAPISpec.Servers, 1)
	assert.Equal(t, "http://localhost", app.openAPISpec.Servers[0].URL)
	assert.Equal(t, "Local server", app.openAPISpec.Servers[0].Description)
}

func TestApp_Description(t *testing.T) {
	app := New(
		SetDescription("New SetDescription"),
	)

	assert.Equal(t, "New SetDescription", app.openAPISpec.Info.Description)
}

func TestApp_Title(t *testing.T) {
	app := New(SetTitle("New SetTitle"))

	assert.Equal(t, "New SetTitle", app.openAPISpec.Info.Title)
}

func TestApp_Version(t *testing.T) {
	app := New(SetVersion("1.0.0"))

	assert.Equal(t, "1.0.0", app.openAPISpec.Info.Version)
}

func TestApp_Contact(t *testing.T) {
	app := New(
		SetContact(&openapi3.Contact{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		}),
	)

	assert.Equal(t, "John Doe", app.openAPISpec.Info.Contact.Name)
	assert.Equal(t, "john.doe@example.com", app.openAPISpec.Info.Contact.Email)
}

func TestApp_License(t *testing.T) {
	app := New(
		SetLicense(&openapi3.License{
			Name: "MIT",
		}),
	)

	assert.Equal(t, "MIT", app.openAPISpec.Info.License.Name)
}

func TestApp_TermsOfService(t *testing.T) {
	app := New(
		SetTermsOfService("https://example.com/terms"),
	)

	assert.Equal(t, "https://example.com/terms", app.openAPISpec.Info.TermsOfService)
}

func TestApp_Setup(t *testing.T) {
	app := New(SetDisableSwagger(true))

	err := app.setup()
	assert.Nil(t, err)
}

func TestApp_Setup_saveOpenAPISpecError(t *testing.T) {
	app := New()

	// Mock yaml.JSONToYAML to simulate an error
	yamlJSONToYAML = mockJSONToYAML
	defer restoreJSONToYAML()

	err := app.setup()
	assert.Error(t, err)
}

func TestOpenAPIHandler(t *testing.T) {
	app := New()

	app.app.Get("/swagger/*", app.openAPIPathHandler)

	t.Run("returns the OpenAPI file", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		_, err := app.app.Test(req)

		assert.NoError(t, err)
	})
}

func TestApp_createDefaultErrorResponses(t *testing.T) {
	app := New()
	responses, err := app.createDefaultErrorResponses()

	assert.Nil(t, err)
	assert.NotNil(t, responses)
	assert.Contains(t, responses, 400) // Assuming 400 is in DefaultErrorResponses
}

func TestApp_createDefaultErrorResponses_error(t *testing.T) {
	app := New()

	realGenerator := generatorNewSchemaRefForValue
	defer func() {
		generatorNewSchemaRefForValue = realGenerator
	}()

	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, errors.New("error")
	}

	_, err := app.createDefaultErrorResponses()
	assert.NotNil(t, err)
}

func TestWriteOpenAPISpec(t *testing.T) {
	jsonData := []byte(`{"openapi": "3.0.3"}`)
	yamlData, err := writeOpenAPISpec(jsonData)

	assert.Nil(t, err)
	assert.NotNil(t, yamlData)
}

// Struct with non-serializable field to simulate error
type NonSerializableStruct struct {
	NonSerializableChan chan int
}

func TestApp_SaveOpenAPISpec_Error(t *testing.T) {
	app := New()

	app.openAPISpec.Components.Schemas["nonSerializable"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Extensions: map[string]interface{}{
				"x-non-serializable": NonSerializableStruct{},
			},
		},
	}

	_, err := app.saveOpenAPISpec()
	assert.NotNil(t, err)
}

func mockJSONToYAML(_ []byte) ([]byte, error) {
	return nil, errors.New("error converting to YAML")
}

var originalJSONToYAML = yaml.JSONToYAML

func restoreJSONToYAML() {
	yamlJSONToYAML = originalJSONToYAML
}

func TestApp_SaveOpenAPISpec_YAMLError(t *testing.T) {
	app := New()

	// Struct that can be serialized to JSON but will fail for YAML conversion
	app.openAPISpec.Components.Schemas["serializable"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Properties: map[string]*openapi3.SchemaRef{
				"valid": {
					Value: &openapi3.Schema{},
				},
			},
		},
	}

	// Mock yaml.JSONToYAML to simulate an error
	yamlJSONToYAML = mockJSONToYAML
	defer restoreJSONToYAML()

	// Attempt to save OpenAPI spec and check for errors
	_, err := app.saveOpenAPISpec()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "error converting to YAML")
}

// Helper function to find an available port
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}

func TestApp_Listen(t *testing.T) {
	app := New()

	// Find a free port to avoid conflicts
	port, err := getFreePort()
	assert.Nil(t, err)

	address := fmt.Sprintf(":%d", port)

	// Run Listen in a separate goroutine since it is blocking
	go func() {
		err = app.Listen(address)
		assert.Nil(t, err)
	}()

	// Wait a bit for the server to start
	// You might want to use a more reliable synchronization mechanism
	// in real tests, like a sync.WaitGroup or a channel.
	<-time.After(time.Second)

	// Attempt to connect to the server
	conn, err := net.Dial("tcp", address)
	assert.Nil(t, err)
	if err == nil {
		conn.Close()
	}

	// Shutdown the server
	assert.NoError(t, app.Shutdown())
}

func TestApp_ListenError(t *testing.T) {
	app := New()

	// Mock yaml.JSONToYAML to simulate an error
	yamlJSONToYAML = mockJSONToYAML
	defer restoreJSONToYAML()

	err := app.Listen(":8080")
	assert.Error(t, err)
}

func TestApp_Run(t *testing.T) {
	// Find a free port to avoid conflicts
	port, err := getFreePort()
	assert.Nil(t, err)

	address := fmt.Sprintf(":%d", port)

	app := New(SetAddress(address))

	// Run Listen in a separate goroutine since it is blocking
	go func() {
		err = app.Run()
		assert.Nil(t, err)
	}()

	// Wait a bit for the server to start
	// You might want to use a more reliable synchronization mechanism
	// in real tests, like a sync.WaitGroup or a channel.
	<-time.After(time.Second)

	// Attempt to connect to the server
	conn, err := net.Dial("tcp", address)
	assert.Nil(t, err)
	if err == nil {
		conn.Close()
	}

	// Shutdown the server
	assert.NoError(t, app.Shutdown())
}

func TestApp_RunError(t *testing.T) {
	app := New()

	// Mock yaml.JSONToYAML to simulate an error
	yamlJSONToYAML = mockJSONToYAML
	defer restoreJSONToYAML()

	err := app.Run()
	assert.Error(t, err)
}

func TestApp_SetAddress(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	New(SetAddress("8080"))
}

func TestApp_SetAddress2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	New(SetAddress(":invalid"))
}

func TestApp_saveOpenAPIToFile(t *testing.T) {
	app := New()

	realOsMkdirAll := osMkdirAll
	defer func() {
		osMkdirAll = realOsMkdirAll
	}()

	osMkdirAll = func(path string, perm os.FileMode) error {
		return assert.AnError
	}

	err := app.saveOpenAPIToFile("test", []byte("test"))
	assert.Error(t, err)
}

func TestApp_saveOpenAPIToFile2(t *testing.T) {
	app := New()

	realOsCreate := osCreate
	defer func() {
		osCreate = realOsCreate
	}()

	osCreate = func(path string) (*os.File, error) {
		return nil, assert.AnError
	}

	err := app.saveOpenAPIToFile("test", []byte("test"))
	assert.Error(t, err)
}

type writeCloserFail struct{}

func (w writeCloserFail) Write(_ []byte) (n int, err error) {
	return 0, assert.AnError
}

func (w writeCloserFail) Close() error {
	return nil
}

func TestApp_saveOpenAPIToFile3(t *testing.T) {
	app := New()

	realWrapperWriteCloser := wrapperWriteCloser
	defer func() {
		wrapperWriteCloser = realWrapperWriteCloser
	}()

	wrapperWriteCloser = func(w io.WriteCloser) io.WriteCloser {
		return &writeCloserFail{}
	}

	err := app.saveOpenAPIToFile("test", []byte("test"))
	assert.Error(t, err)
}

func TestApp_Setup_saveOpenAPIToFileError(t *testing.T) {
	realOsMkdirAll := osMkdirAll
	defer func() {
		osMkdirAll = realOsMkdirAll
	}()

	osMkdirAll = func(path string, perm os.FileMode) error {
		return assert.AnError
	}

	// Find a free port to avoid conflicts
	port, err := getFreePort()
	assert.Nil(t, err)

	address := fmt.Sprintf(":%d", port)

	app := New(SetAddress(address))

	// Run Listen in a separate goroutine since it is blocking
	go func() {
		err = app.Run()
		Use(app, fiberrecover.New())
		assert.Nil(t, err)
	}()

	// Wait a bit for the server to start
	// You might want to use a more reliable synchronization mechanism
	// in real tests, like a sync.WaitGroup or a channel.
	<-time.After(time.Second)

	// Attempt to connect to the server
	conn, err := net.Dial("tcp", address)
	assert.Nil(t, err)
	if err == nil {
		conn.Close()
	}

	// Shutdown the server
	assert.NoError(t, app.Shutdown())
}
