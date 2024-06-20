package lite

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"reflect"
	"testing"
)

type DeserializerTestSuite struct {
	suite.Suite
}

func (suite *DeserializerTestSuite) SetupTest() {
}

func TestDeserializerTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(DeserializerTestSuite))
}

func (suite *DeserializerTestSuite) TestDeserialize() {
	type testStruct struct {
		Authorization string `lite:""`
	}

	var test = testStruct{}

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(nil, val, nil)
	assert.Error(suite.T(), err)
}

type multipartForm struct {
	Body Image `lite:"req=body,multipart/form-data"`
}

type Image struct {
	Image *multipart.FileHeader `form:"image"`
	Name  string                `form:"name"`
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyMultipart() {
	var test = multipartForm{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[multipartForm](ctx, app, "/foo")
	c.Request().Header.SetContentType("multipart/form-data" + `;boundary="b"`)

	// fake image PNG in base64
	imageBase64 := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/wcAAgEBABkAB4cAAAAASUVORK5CYII="
	imageBytes, _ := base64.StdEncoding.DecodeString(imageBase64)

	var body []byte
	body = append(body, []byte("--b\r\nContent-Disposition: form-data; name=\"name\"\r\n\r\njohn\r\n--b\r\nContent-Disposition: form-data; name=\"image\"; filename=\"fakeimage.png\"\r\nContent-Type: image/png\r\n\r\n")...)
	body = append(body, imageBytes...)
	body = append(body, []byte("\r\n--b--")...)

	c.Request().SetBody(body)
	c.Request().Header.SetContentLength(len(body))

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.NoError(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyMultipartError() {
	var test = multipartForm{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[multipartForm](ctx, app, "/foo")
	defer c.Request().Reset()

	c.Request().Header.SetContentType("multipart/form-data" + `;boundary="b"`)

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.Error(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyApplicationOctetStream() {
	type testStruct struct {
		Body []byte `lite:"req=body,application/octet-stream"`
	}

	var test = testStruct{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[testStruct](ctx, app, "/foo")
	c.Request().Header.SetContentType("application/octet-stream")

	// fake image PNG in base64
	imageBase64 := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/wcAAgEBABkAB4cAAAAASUVORK5CYII="
	imageBytes, _ := base64.StdEncoding.DecodeString(imageBase64)

	c.Request().SetBody(imageBytes)
	c.Request().Header.SetContentLength(len(imageBytes))

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.NoError(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyApplicationOctetStreamError() {
	var test = multipartForm{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[multipartForm](ctx, app, "/foo")
	defer c.Request().Reset()

	c.Request().Header.SetContentType("application/octet-stream")

	// fake image PNG in base64
	imageBase64 := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/wcAAgEBABkAB4cAAAAASUVORK5CYII="
	imageBytes, _ := base64.StdEncoding.DecodeString(imageBase64)

	c.Request().SetBody(imageBytes)
	c.Request().Header.SetContentLength(len(imageBytes))

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.Error(suite.T(), err)
}

type Val struct {
	Name string `form:"name"`
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyApplicationXWwwFormUrlencoded() {
	type testStruct struct {
		Body Val `lite:"req=body,application/x-www-form-urlencoded"`
	}

	var test = testStruct{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[testStruct](ctx, app, "/foo")
	c.Request().Header.SetContentType("application/x-www-form-urlencoded")

	c.Request().SetBody([]byte("name=john"))
	c.Request().Header.SetContentLength(9)

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.NoError(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyTextHTML() {
	type testStruct struct {
		Body string `lite:"req=body,text/html"`
	}

	var test = testStruct{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[testStruct](ctx, app, "/foo")
	c.Request().Header.SetContentType("text/html")

	c.Request().SetBody([]byte("test"))
	c.Request().Header.SetContentLength(4)

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.NoError(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyApplicationPDF() {
	type testStruct struct {
		Body []byte `lite:"req=body,application/pdf"`
	}

	var test = testStruct{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[testStruct](ctx, app, "/foo")
	c.Request().Header.SetContentType("application/pdf")

	c.Request().SetBody([]byte("test"))
	c.Request().Header.SetContentLength(4)

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.NoError(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyApplicationZip() {
	type testStruct struct {
		Body []byte `lite:"req=body,application/zip"`
	}

	var test = testStruct{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[testStruct](ctx, app, "/foo")
	c.Request().Header.SetContentType("application/zip")

	c.Request().SetBody([]byte("test"))
	c.Request().Header.SetContentLength(4)

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.NoError(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyImage() {
	type testStruct struct {
		Body []byte `lite:"req=body,image/*"`
	}

	var test = testStruct{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[testStruct](ctx, app, "/foo")
	c.Request().Header.SetContentType("image/png")

	c.Request().SetBody([]byte("test"))
	c.Request().Header.SetContentLength(4)

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.NoError(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyUnsupportedContentType() {
	type testStruct struct {
		Body []byte `lite:"req=body,application/test"`
	}

	var test = testStruct{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[testStruct](ctx, app, "/foo")
	c.Request().Header.SetContentType("application/test")

	c.Request().SetBody([]byte("test"))
	c.Request().Header.SetContentLength(4)

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.Error(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestSetFieldValue() {
	type testStruct struct {
		Name    string  `form:"name"`
		Age     int     `form:"age"`
		Price   float64 `form:"price"`
		IsAdmin bool    `form:"is_admin"`
		Int8    int8    `form:"int8"`
		Int16   int16   `form:"int16"`
		Int32   int32   `form:"int32"`
		Int64   int64   `form:"int64"`
		Uint    uint    `form:"uint"`
		Uint8   uint8   `form:"uint8"`
		Uint16  uint16  `form:"uint16"`
		Uint32  uint32  `form:"uint32"`
		Uint64  uint64  `form:"uint64"`
		F32     float32 `form:"f32"`
		F64     float64 `form:"f64"`
		Slice   []int   `form:"slice"`
		Fun     func()  `form:"fun"`
		Val     Val     `form:"val"`
		Byt     []byte  `form:"byt"`
	}

	var test = testStruct{}

	val := reflect.ValueOf(&test).Elem()

	err := setFieldValue(val.Field(0), "john")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(1), "1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(1), "true")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(2), "1.1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(2), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(3), "true")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(3), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(4), "1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(4), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(5), "1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(5), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(6), "1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(6), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(7), "1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(7), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(8), "1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(8), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(9), "1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(9), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(10), "1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(10), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(11), "1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(11), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(12), "1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(12), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(13), "1.1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(13), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(14), "1.1")
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(14), "test")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(15), "1")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(16), "1,2,3")
	assert.Error(suite.T(), err)

	err = setFieldValue(val.Field(17), `{"name":"john"}`)
	assert.NoError(suite.T(), err)

	err = setFieldValue(val.Field(18), "test")
	assert.NoError(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestMapToStruct() {
	type testStruct struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	m := map[string]any{
		"name": "john",
		"age":  "1",
	}

	val := reflect.ValueOf(&testStruct{}).Elem()

	err := mapToStruct(m, val.Addr().Interface())
	assert.NoError(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestMapToStructError() {
	type testStruct struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	m := map[string]any{
		"name": "john",
		"age":  "test",
	}

	val := reflect.ValueOf(&testStruct{}).Elem()

	err := mapToStruct(m, val.Addr().Interface())
	assert.Error(suite.T(), err)
}

func (suite *DeserializerTestSuite) TestMapToStructWithEmptyTag() {
	type testStruct struct {
		Name string `form:""`
		Age  int    `form:"age"`
	}

	m := map[string]any{
		"name": "john",
		"age":  "1",
	}

	val := reflect.ValueOf(&testStruct{}).Elem()

	err := mapToStruct(m, val.Addr().Interface())
	assert.NoError(suite.T(), err)
}
