package lite

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"reflect"
	"testing"
	"time"
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
		Name string `lite:"params=name"`
	}

	type TestResponse struct {
		Name string `json:"name"`
	}

	_, err = registerOpenAPIOperation[TestResponse, TestRequest](app, "GET", "/test", "application/json", 200)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestRegisterOpenAPIOperationResponseNil(t *testing.T) {
	app := New()

	var err error

	type TestRequest struct {
		Name string `lite:"params=name"`
	}

	_, err = registerOpenAPIOperation[any, TestRequest](app, "GET", "/test", "application/json", 200)
	if err != nil {
		t.Fatal("should be error")
	}
}

func TestRegisterOpenAPIOperationGenerateBodyStringError(t *testing.T) {
	app := New()

	var err error

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	_, err = registerOpenAPIOperation[string, string](app, "GET", "/test", "application/json", 200)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestRegisterOpenAPIOperationGenerateSliceByteError(t *testing.T) {
	app := New()

	var err error

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	_, err = registerOpenAPIOperation[[]byte, []byte](app, "GET", "/test", "application/json", 200)
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

func TestRegisterOpenAPIOperationGenerateError4(t *testing.T) {
	app := New()

	var err error

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue

	type TestRequest struct{}

	type TestResponse struct {
		FirstName string `json:"first_name"`
	}

	_, err = registerOpenAPIOperation[TestResponse, TestRequest](app, "GET", "/test", "application/json", 200)
	if err != nil {
		t.Fatal(err)
	}

	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	type TestResponse2 struct {
		LastName *string `json:"last_name"`
	}

	_, err = registerOpenAPIOperation[TestResponse2, TestRequest](app, "GET", "/test", "application/json", 200)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestRegisterOpenAPIOperationGenerateError5(t *testing.T) {
	app := New()

	var err error

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue

	type TestRequest struct{}

	_, err = registerOpenAPIOperation[string, TestRequest](app, "GET", "/test", "application/json", 200)
	if err != nil {
		t.Fatal(err)
	}

	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	_, err = registerOpenAPIOperation[*string, TestRequest](app, "GET", "/test", "application/json", 200)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestRegisterOpenAPIOperation(t *testing.T) {
	app := New()

	var err error

	type TestRequest struct{}

	_, err = registerOpenAPIOperation[string, TestRequest](app, "GET", "/test", "application/json", 200)
	if err != nil {
		t.Fatal(err)
	}

	_, err = registerOpenAPIOperation[*string, TestRequest](app, "GET", "/test", "application/json", 200)
	if err != nil {
		t.Fatal("should not be error")
	}
}

type testStruct struct {
	Name    string    `json:"name" xml:"name" form:"name"`
	Age     int       `json:"age" xml:"age" form:"age"`
	Address string    `json:"address" xml:"address" form:"Address"`
	Time    time.Time `json:"time" xml:"time" form:"time"`
}

func TestGetRequiredValue_Struct(t *testing.T) {
	contentType := "application/json"
	fieldType := reflect.TypeOf(testStruct{})
	schema := &openapi3.Schema{
		Properties: map[string]*openapi3.SchemaRef{
			"name":    {Value: &openapi3.Schema{}},
			"age":     {Value: &openapi3.Schema{}},
			"address": {Value: &openapi3.Schema{}},
			"time":    {Value: &openapi3.Schema{}},
		},
	}

	ok := getRequiredValue(contentType, fieldType, schema)
	if !ok {
		t.Errorf("expected true, got false")
	}
	if len(schema.Required) != 4 {
		t.Errorf("expected 4 required fields, got %d", len(schema.Required))
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
			"time":    {Value: &openapi3.Schema{}},
		},
	}

	ok := getRequiredValue(contentType, fieldType, schema)
	if !ok {
		t.Errorf("expected true, got false")
	}
	if len(schema.Required) != 4 {
		t.Errorf("expected 4 required fields, got %d", len(schema.Required))
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
			"time":    {Value: &openapi3.Schema{}},
		},
	}

	ok := getRequiredValue(contentType, fieldType, schema)
	if !ok {
		t.Errorf("expected true, got false")
	}
	if len(schema.Required) != 4 {
		t.Errorf("expected 4 required fields, got %d", len(schema.Required))
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
	if !ok {
		t.Errorf("expected false, got true")
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
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterPtrInt2(t *testing.T) {
	testErr := testPtrInt{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err != nil {
		t.Fatal(err)
	}
}

type CreateBody struct {
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type Params struct {
	ID uint64 `lite:"params=id"`
}

type testReq struct {
	Authorization *string `lite:"header=Authorization,isauth,scheme=bearer,name=Authorization"`
	Name          string  `lite:"header=name"`
	Params        Params
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
	Params        Params
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
	query query
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

type testReq6 struct {
	Body *multipart.FileHeader `lite:"req=body"`
}

func TestRegisterBodyMultipartFileHeader(t *testing.T) {
	testErr := testReq6{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err != nil {
		t.Fatal(err)
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
	if tag != "unknown" {
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
	ID      uint64 `lite:"params=id"`
	IsAdmin string `lite:"params=is_admin"`
}

type metadata struct {
	FirstName string    `form:"first_name" json:"first_name"`
	LastName  string    `form:"last_name" json:"last_name"`
	Birthday  time.Time `form:"birthday" json:"birthday"`
}

type bodyRequest struct {
	Name     string                `form:"name"`
	File     *multipart.FileHeader `form:"file"`
	Metadata *metadata             `form:"metadata"`
}

type testRequest struct {
	Params params
	Filter *string      `lite:"query=filter"`
	Age    string       `lite:"query=age"`
	Token  string       `lite:"header=token"`
	Value  *string      `lite:"header=value"`
	Cookie *http.Cookie `lite:"cookie=cookie"`
	Body   bodyRequest  `lite:"req=body,multipart/form-data"`
}

type bodyMultiFile struct {
	Files []*multipart.FileHeader `form:"files"`
}

type testRequestMultiFileRequest struct {
	Body bodyMultiFile `lite:"req=body,multipart/form-data"`
}

type testResponse struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Gender      Gender   `json:"gender" enums:"male,female"`
	GenderSlice []Gender `json:"gender_slice" enums:"male,female"`
}

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

func TestRegisterSetParamSchemaError(t *testing.T) {
	type testErrorParams struct {
		ID  uint64 `lite:"params=id"`
		Age string `lite:"params=age"`
	}

	testErr := testErrorParams{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err != nil {
		t.Fatal("should be error")
	}

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	err = register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestRegisterSetHeaderSchemeError(t *testing.T) {
	type testErrorHeader struct {
		ID  uint64 `lite:"header=id"`
		Age string `lite:"header=age"`
	}

	testErr := testErrorHeader{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err != nil {
		t.Fatal("should be error")
	}

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	err = register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

func TestRegisterSetSecurityScheme(t *testing.T) {
	type testHeader struct {
		BasicAuth string `lite:"header=Authorization,isauth,scheme=basic"`
	}

	testErr := testHeader{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterSetSecurityScheme2(t *testing.T) {
	type testHeader struct {
		BearerAuth string `lite:"header=Authorization,isauth,scheme=bearer,type=test"`
	}

	testErr := testHeader{}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testErr)

	err := register(app, operation, dstVal)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterSetBodySchema(t *testing.T) {
	type testStruct struct {
		body string `lite:"req=body"`
	}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testStruct{})

	err := register(app, operation, dstVal)
	if err != nil {
		t.Fatal(err)
	}

	type testStruct2 struct {
		body *string `lite:"req=body"`
	}

	dstVal = reflect.ValueOf(testStruct2{})

	err = register(app, operation, dstVal)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterSetBodySchemaError(t *testing.T) {
	type testStruct struct {
		body string `lite:"req=body"`
	}

	app := New()
	operation := openapi3.NewOperation()
	dstVal := reflect.ValueOf(testStruct{})

	err := register(app, operation, dstVal)
	if err != nil {
		t.Fatal(err)
	}

	realgeneratorNewSchemaRefForValue := generatorNewSchemaRefForValue
	generatorNewSchemaRefForValue = func(value interface{}, schemas openapi3.Schemas) (*openapi3.SchemaRef, error) {
		return nil, assert.AnError
	}
	defer func() {
		generatorNewSchemaRefForValue = realgeneratorNewSchemaRefForValue
	}()

	err = register(app, operation, dstVal)
	if err == nil {
		t.Fatal("should be error")
	}
}

type MyType struct{}
type MyParamType struct{}
type AnotherParamType struct{}

func TestExtractBaseName(t *testing.T) {
	tests := []struct {
		name      string
		typeInput reflect.Type
		expected  string
	}{
		{
			name:      "generic type with slash",
			typeInput: reflect.TypeOf(map[MyType]AnotherParamType{}),
			expected:  "mapMyType",
		},
		{
			name:      "generic type with dot",
			typeInput: reflect.TypeOf(map[MyType]MyParamType{}),
			expected:  "mapMyType",
		},
		{
			name:      "generic type without slash or dot",
			typeInput: reflect.TypeOf(map[MyType]MyParamType{}),
			expected:  "mapMyType",
		},
		{
			name:      "non-generic type",
			typeInput: reflect.TypeOf(MyType{}),
			expected:  "MyType",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractBaseName(tt.typeInput)
			assert.Equal(t, tt.expected, result)
		})
	}
}
