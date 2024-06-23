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
	spec := NewOpenAPISpec()

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
	app := New()

	assert.NotNil(t, app.App)
	assert.Equal(t, "OpenAPI", app.OpenAPISpec.Info.Title)
	assert.Equal(t, defaultOpenAPIConfig.SwaggerURL, app.OpenAPIConfig.SwaggerURL)
	assert.Equal(t, defaultOpenAPIConfig.YamlURL, app.OpenAPIConfig.YamlURL)
}

func TestApp_AddTags(t *testing.T) {
	app := New()
	app.AddTags("tag1", "tag2")

	assert.Contains(t, app.tags, "tag1")
	assert.Contains(t, app.tags, "tag2")
}

func TestApp_SaveOpenAPISpec(t *testing.T) {
	app := New()
	yamlData, err := app.SaveOpenAPISpec()

	assert.Nil(t, err)
	assert.NotNil(t, yamlData)
}

func TestApp_AddServer(t *testing.T) {
	app := New()
	app.AddServer("http://localhost", "Local server")

	assert.Len(t, app.OpenAPISpec.Servers, 1)
	assert.Equal(t, "http://localhost", app.OpenAPISpec.Servers[0].URL)
	assert.Equal(t, "Local server", app.OpenAPISpec.Servers[0].Description)
}

func TestApp_Description(t *testing.T) {
	app := New()
	app.Description("New Description")

	assert.Equal(t, "New Description", app.OpenAPISpec.Info.Description)
}

func TestApp_Title(t *testing.T) {
	app := New()
	app.Title("New Title")

	assert.Equal(t, "New Title", app.OpenAPISpec.Info.Title)
}

func TestApp_Version(t *testing.T) {
	app := New()
	app.Version("1.0.0")

	assert.Equal(t, "1.0.0", app.OpenAPISpec.Info.Version)
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

	app.OpenAPISpec.Components.Schemas["nonSerializable"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Extensions: map[string]interface{}{
				"x-non-serializable": NonSerializableStruct{},
			},
		},
	}

	_, err := app.SaveOpenAPISpec()
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
	app.OpenAPISpec.Components.Schemas["serializable"] = &openapi3.SchemaRef{
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
	_, err := app.SaveOpenAPISpec()

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
