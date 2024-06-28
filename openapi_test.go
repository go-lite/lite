package lite

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"reflect"
	"testing"
)

func TestRegisterOpenAPIOperationGenerateError(t *testing.T) {
	app := New()

	var err error

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	type TestRequest struct {
		Name string `lite:"path=name"`
	}

	type TestResponse struct {
		Name string `json:"name"`
	}

	_, err = registerOpenAPIOperation[TestResponse, TestRequest](app, "GET", "/test", "application/json", 200)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestRegisterOpenAPIOperationGenerateError2(t *testing.T) {
	app := New()

	var err error

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	type TestRequest struct{}

	type TestResponse struct {
		Name string `json:"name"`
	}

	_, err = registerOpenAPIOperation[TestResponse, TestRequest](app, "GET", "/test", "application/json", 200)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestRegisterOpenAPIOperationGenerateError3(t *testing.T) {
	app := New()

	var err error

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	type TestRequest struct{}

	type TestResponse struct{}

	_, err = registerOpenAPIOperation[TestResponse, TestRequest](app, "GET", "/test", "application/json", 200)
	if err == nil {
		t.Fatal("should be error")
	}
}

type testStruct struct {
	Name    string `json:"name" xml:"name" form:"name"`
	Age     int    `json:"age" xml:"age" form:"age"`
	Address string `json:"address" xml:"address" form:"Address"`
}

func TestGetRequiredValue_Struct(t *testing.T) {
	contentType := "application/json"
	fieldType := reflect.TypeOf(testStruct{})
	schema := &openapi3.Schema{
		Properties: map[string]*openapi3.SchemaRef{
			"name":    {Value: &openapi3.Schema{}},
			"age":     {Value: &openapi3.Schema{}},
			"address": {Value: &openapi3.Schema{}},
		},
	}

	ok := getRequiredValue(contentType, fieldType, schema)
	if !ok {
		t.Errorf("expected true, got false")
	}
	if len(schema.Required) != 3 {
		t.Errorf("expected 3 required fields, got %d", len(schema.Required))
	}
}

func TestGetRequiredValue_StructWithXML(t *testing.T) {
	contentType := "application/xml"
	fieldType := reflect.TypeOf(testStruct{})
	schema := &openapi3.Schema{
		Properties: map[string]*openapi3.SchemaRef{
			"name":    {Value: &openapi3.Schema{}},
			"age":     {Value: &openapi3.Schema{}},
			"address": {Value: &openapi3.Schema{}},
		},
	}

	ok := getRequiredValue(contentType, fieldType, schema)
	if !ok {
		t.Errorf("expected true, got false")
	}
	if len(schema.Required) != 3 {
		t.Errorf("expected 3 required fields, got %d", len(schema.Required))
	}
}

func TestGetRequiredValue_StructWithForm(t *testing.T) {
	contentType := "application/x-www-form-urlencoded"
	fieldType := reflect.TypeOf(testStruct{})
	schema := &openapi3.Schema{
		Properties: map[string]*openapi3.SchemaRef{
			"name":    {Value: &openapi3.Schema{}},
			"age":     {Value: &openapi3.Schema{}},
			"address": {Value: &openapi3.Schema{}},
		},
	}

	ok := getRequiredValue(contentType, fieldType, schema)
	if !ok {
		t.Errorf("expected true, got false")
	}
	if len(schema.Required) != 3 {
		t.Errorf("expected 3 required fields, got %d", len(schema.Required))
	}
}

func TestGetRequiredValue_Slice(t *testing.T) {
	contentType := "application/json"
	fieldType := reflect.TypeOf([]int{})
	schema := &openapi3.Schema{
		Items: &openapi3.SchemaRef{
			Value: &openapi3.Schema{},
		},
	}

	ok := getRequiredValue(contentType, fieldType, schema)
	if !ok {
		t.Errorf("expected true, got false")
	}
}

func TestGetRequiredValue_SliceOfStructWithPointer(t *testing.T) {
	contentType := "application/json"
	fieldType := reflect.TypeOf([]*testStruct{})
	schema := &openapi3.Schema{
		Items: &openapi3.SchemaRef{
			Value: &openapi3.Schema{},
		},
	}

	ok := getRequiredValue(contentType, fieldType, schema)
	if ok {
		t.Errorf("expected true, got false")
	}
}

func TestGetRequiredValue_SliceOfUint8(t *testing.T) {
	contentType := "application/json"
	fieldType := reflect.TypeOf([]byte{})
	schema := &openapi3.Schema{
		Items: &openapi3.SchemaRef{
			Value: &openapi3.Schema{},
		},
	}

	ok := getRequiredValue(contentType, fieldType, schema)
	if !ok {
		t.Errorf("expected true, got false")
	}
}

func TestGetRequiredValue_Interface(t *testing.T) {
	contentType := "application/json"
	fieldType := reflect.TypeOf((*interface{})(nil)).Elem()
	schema := &openapi3.Schema{}

	ok := getRequiredValue(contentType, fieldType, schema)
	if ok {
		t.Errorf("expected false, got true")
	}
}

func TestGetRequiredValue_Pointer(t *testing.T) {
	contentType := "application/json"
	fieldType := reflect.TypeOf((*int)(nil))
	schema := &openapi3.Schema{}

	ok := getRequiredValue(contentType, fieldType, schema)
	if ok {
		t.Errorf("expected false, got true")
	}
}

func TestGetRequiredValue_PointerOfStruct(t *testing.T) {
	contentType := "application/xml"
	fieldType := reflect.TypeOf((*testStruct)(nil))
	schema := &openapi3.Schema{}

	ok := getRequiredValue(contentType, fieldType, schema)
	if ok {
		t.Errorf("expected false, got true")
	}
}

func TestGetRequiredValue_Map(t *testing.T) {
	contentType := "application/json"
	fieldType := reflect.TypeOf(map[string]int{})
	schema := &openapi3.Schema{
		AdditionalProperties: openapi3.AdditionalProperties{
			Schema: &openapi3.SchemaRef{},
		},
	}

	ok := getRequiredValue(contentType, fieldType, schema)
	if ok {
		t.Errorf("expected true, got false")
	}
}

// this test will panic
func TestGetRequiredValue_ChanPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	contentType := "application/json"
	fieldType := reflect.TypeOf(make(chan int))
	schema := &openapi3.Schema{}

	getRequiredValue(contentType, fieldType, schema)
}

func TestGetRequiredValue_BasicTypes(t *testing.T) {
	contentType := "application/json"
	basicTypes := []reflect.Type{
		reflect.TypeOf(true),
		reflect.TypeOf(1),
		reflect.TypeOf(1.0),
		reflect.TypeOf("string"),
	}

	for _, fieldType := range basicTypes {
		schema := &openapi3.Schema{}

		ok := getRequiredValue(contentType, fieldType, schema)
		if !ok {
			t.Errorf("expected true, got false for type %v", fieldType)
		}
	}
}

type testError struct {
	Ch chan int `lite:"req=body"`
}

type testError2 struct {
	In **int `lite:"req=body"`
}

func TestRegisterChan(t *testing.T) {
	testErr := testError{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestRegisterPtrPtrInt(t *testing.T) {
	var val *int
	val = new(int)
	testErr := testError2{
		In: &val,
	}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

type testErrorStruct struct {
	In *testError2 `lite:"req=body"`
}

func TestRegisterPtrStruct(t *testing.T) {
	var val *testError2
	val = new(testError2)
	testErr := testErrorStruct{
		In: val,
	}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

type testPtrInt struct {
	In *int `lite:"req=body"`
}

func TestRegisterPtrInt(t *testing.T) {
	var val *int
	val = new(int)
	testErr := testPtrInt{
		In: val,
	}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestRegisterPtrInt2(t *testing.T) {
	testErr := testPtrInt{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

type CreateBody struct {
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type Params struct {
	ID uint64 `lite:"path=id"`
}

type testReq struct {
	Authorization *string `lite:"header=Authorization,isauth,scheme=bearer,name=Authorization"`
	Name          string  `lite:"header=name"`
	ID            Params
	Body          CreateBody `lite:"req=body,application/json"`
}

func TestRegisterStruct(t *testing.T) {
	testErr := testReq{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterStructGeneratorNewSchemaRefForValueError(t *testing.T) {
	testErr := testReq{}

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

type testReq2 struct {
	Authorization *string `lite:"header=Authorization,isauth,scheme=bearer"`
	ID            Params
	Body          CreateBody `lite:"req=body,application/xml,application/json"`
}

// this test will panic
func TestRegisterStructPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	testErr := testReq2{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

type query struct {
	ID   uint64   `lite:"query=id"`
	Name **string `lite:"query=name"`
}

type testReq3 struct {
	ID query
}

func TestRegisterQueryError(t *testing.T) {
	testErr := testReq3{}

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

type body struct {
	Name string                `form:"name"`
	File *multipart.FileHeader `form:"name"`
}

type testReq4 struct {
	Body body `lite:"req=body"`
}

func TestRegisterBodyError(t *testing.T) {
	testErr := testReq4{}

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

type testReq5 struct {
	Cookie *http.Cookie `lite:"cookie=cookie"`
}

func TestRegisterCookie(t *testing.T) {
	testErr := testReq5{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterCookieGeneratorNewSchemaRefForValueError(t *testing.T) {
	testErr := testReq5{}

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestGetStructTag(t *testing.T) {
	contentType := "text/plain"

	tag := getStructTag(contentType)
	if tag != "txt" {
		t.Errorf("expected txt, got %s", tag)
	}

	contentType = "application/json"

	tag = getStructTag(contentType)
	if tag != "json" {
		t.Errorf("expected json, got %s", tag)
	}

	contentType = "application/xml"

	tag = getStructTag(contentType)
	if tag != "xml" {
		t.Errorf("expected xml, got %s", tag)
	}

	contentType = "application/x-www-form-urlencoded"

	tag = getStructTag(contentType)
	if tag != "form" {
		t.Errorf("expected form, got %s", tag)
	}

	contentType = "multipart/form-data"

	tag = getStructTag(contentType)
	if tag != "form" {
		t.Errorf("expected form, got %s", tag)
	}

	contentType = "text/plain"

	tag = getStructTag(contentType)
	if tag != "txt" {
		t.Errorf("expected txt, got %s", tag)
	}

	contentType = "application/octet-stream"

	tag = getStructTag(contentType)
	if tag != "binary" {
		t.Errorf("expected binary, got %s", tag)
	}

	contentType = "application/pdf"

	tag = getStructTag(contentType)
	if tag != "pdf" {
		t.Errorf("expected pdf, got %s", tag)
	}

	contentType = "image/png"

	tag = getStructTag(contentType)
	if tag != "png" {
		t.Errorf("expected png, got %s", tag)
	}

	contentType = "image/jpeg"

	tag = getStructTag(contentType)
	if tag != "jpeg" {
		t.Errorf("expected jpeg, got %s", tag)
	}

	contentType = "application/fake+json"

	tag = getStructTag(contentType)
	if tag != "json" {
		t.Errorf("expected json, got %s", tag)
	}
}

func TestTagFromType(t *testing.T) {
	var v interface{}

	tag := tagFromType(v)
	if tag != "unknown-interface" {
		t.Errorf("expected unknown-interface, got %s", tag)
	}

	v = new(int)

	tag = tagFromType(v)
	if tag != "int" {
		t.Errorf("expected int, got %s", tag)
	}

	v = new(string)

	tag = tagFromType(v)
	if tag != "string" {
		t.Errorf("expected string, got %s", tag)
	}

	v = new(bool)

	tag = tagFromType(v)
	if tag != "bool" {
		t.Errorf("expected bool, got %s", tag)
	}

	v = new(float64)

	tag = tagFromType(v)
	if tag != "float64" {
		t.Errorf("expected float64, got %s", tag)
	}

	v = new(float32)

	tag = tagFromType(v)
	if tag != "float32" {
		t.Errorf("expected float32, got %s", tag)
	}

	v = new(complex64)

	tag = tagFromType(v)
	if tag != "complex64" {
		t.Errorf("expected complex64, got %s", tag)
	}
}

func TestDive(t *testing.T) {
	var v []int

	val := dive(reflect.TypeOf(v), 0)
	if val != "default" {
		t.Errorf("expected default, got %s", val)
	}
}

type params struct {
	ID      uint64 `lite:"path=id"`
	IsAdmin string `lite:"path=is_admin"`
}

type metadata struct {
	FirstName string `form:"first_name" json:"first_name"`
	LastName  string `form:"last_name" json:"last_name"`
}

type bodyRequest struct {
	Name     string                `form:"name"`
	File     *multipart.FileHeader `form:"file"`
	Metadata *metadata             `form:"metadata"`
}

type testRequest struct {
	Params params
	Filter *string      `lite:"query=filter"`
	Cookie *http.Cookie `lite:"cookie=cookie"`
	Body   bodyRequest  `lite:"req=body,multipart/form-data"`
}

type testResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func TestOpenAPI(t *testing.T) {
	app := New()

	_, err := registerOpenAPIOperation[testResponse, testRequest](app, "POST", "/test/:id/:is_admin", "application/json", 200)
	if err != nil {
		t.Fatal(err)
	}

	spec, err := app.SaveOpenAPISpec()
	if err != nil {
		t.Fatal(err)
	}

	expected := `components:
    parameters:
        cookie:
            in: cookie
            name: cookie
            schema:
                $ref: '#/components/schemas/cookie'
        filter:
            in: query
            name: filter
            schema:
                $ref: '#/components/schemas/filter'
        id:
            in: path
            name: id
            required: true
            schema:
                $ref: '#/components/schemas/id'
        is_admin:
            in: path
            name: is_admin
            required: true
            schema:
                $ref: '#/components/schemas/is_admin'
    schemas:
        bodyRequest:
            properties:
                file:
                    format: byte
                    type: string
                metadata:
                    properties:
                        first_name:
                            type: string
                        last_name:
                            type: string
                    type: object
                name:
                    type: string
            required:
                - name
                - file
            type: object
        cookie:
            properties:
                Domain:
                    type: string
                Expires:
                    format: date-time
                    type: string
                HttpOnly:
                    type: boolean
                MaxAge:
                    type: integer
                Name:
                    type: string
                Path:
                    type: string
                Raw:
                    type: string
                RawExpires:
                    type: string
                SameSite:
                    type: integer
                Secure:
                    type: boolean
                Unparsed:
                    items:
                        type: string
                    type: array
                Value:
                    type: string
            type: object
        filter:
            type: string
        httpGenericError:
            properties:
                id:
                    type: string
                message:
                    type: string
                status:
                    type: integer
            type: object
        id:
            maximum: 1.8446744073709552e+19
            minimum: 0
            type: integer
        is_admin:
            type: string
        testResponse:
            properties:
                first_name:
                    type: string
                id:
                    maximum: 1.8446744073709552e+19
                    minimum: 0
                    type: integer
                last_name:
                    type: string
                name:
                    type: string
            required:
                - id
                - name
                - first_name
                - last_name
            type: object
info:
    description: OpenAPI
    title: OpenAPI
    version: 0.0.1
openapi: 3.0.3
paths:
    /test/{id}/{is_admin}:
        post:
            operationId: POST/test/:id/:is_admin
            parameters:
                - $ref: '#/components/parameters/id'
                - $ref: '#/components/parameters/is_admin'
                - $ref: '#/components/parameters/filter'
                - $ref: '#/components/parameters/cookie'
            requestBody:
                content:
                    multipart/form-data:
                        schema:
                            $ref: '#/components/schemas/bodyRequest'
            responses:
                "200":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/testResponse'
                    description: OK
                "400":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Bad Request
                "401":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Unauthorized
                "404":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Not Found
                "409":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Conflict
                "500":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Internal Server Error`

	assert.YAMLEqf(t, expected, string(spec), "openapi generated spec")
}
