package lite

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/yaml"
	"net"
	"testing"
	"time"

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
		SetSwaggerURL("/swagger/*"),
		SetOpenAPIPath("/api/openapi.json"),
		SetUIHandler(defaultOpenAPIHandler),
		SetDisableLocalSave(false),
		SetTypeOfExtension(JSONExtension),
		SetAddress(":8080"),
	)

	assert.NotNil(t, app.app)
	assert.Equal(t, "OpenAPI", app.openAPISpec.Info.Title)
	assert.Equal(t, defaultOpenAPIConfig.swaggerURL, app.openAPIConfig.swaggerURL)
	assert.Equal(t, "/api/openapi.json", app.openAPIConfig.openapiPath)
}

func TestApp_AddTags(t *testing.T) {
	app := New()
	app.AddTags("tag1", "tag2")

	assert.Contains(t, app.tags, "tag1")
	assert.Contains(t, app.tags, "tag2")
}

func TestApp_SaveOpenAPISpec(t *testing.T) {
	app := New()
	yamlData, err := app.saveOpenAPISpec()

	assert.Nil(t, err)
	assert.NotNil(t, yamlData)
}

func TestApp_AddServer(t *testing.T) {
	app := New()
	app.AddServer("http://localhost", "Local server")

	assert.Len(t, app.openAPISpec.Servers, 1)
	assert.Equal(t, "http://localhost", app.openAPISpec.Servers[0].URL)
	assert.Equal(t, "Local server", app.openAPISpec.Servers[0].Description)
}

func TestApp_Description(t *testing.T) {
	app := New()
	app.Description("New Description")

	assert.Equal(t, "New Description", app.openAPISpec.Info.Description)
}

func TestApp_Title(t *testing.T) {
	app := New()
	app.Title("New Title")

	assert.Equal(t, "New Title", app.openAPISpec.Info.Title)
}

func TestApp_Version(t *testing.T) {
	app := New()
	app.Version("1.0.0")

	assert.Equal(t, "1.0.0", app.openAPISpec.Info.Version)
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

func mockJSONToYAML(d []byte) ([]byte, error) {
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
		err := app.Listen(address)
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

func TestApp_Run(t *testing.T) {
	// Find a free port to avoid conflicts
	port, err := getFreePort()
	assert.Nil(t, err)

	address := fmt.Sprintf(":%d", port)

	app := New(SetAddress(address))

	// Run Listen in a separate goroutine since it is blocking
	go func() {
		err := app.Run()
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
