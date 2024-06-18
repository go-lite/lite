package lite

import (
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
}

func (suite *DeserializerTestSuite) TestDeserializeWithBodyMultipart() {
	var test = multipartForm{}

	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[multipartForm](ctx, app, "/foo")
	c.Request().Header.SetContentType("multipart/form-data" + `;boundary="b"`)
	body := []byte("--b\r\nContent-Disposition: form-data; name=\"name\"\r\n\r\njohn\r\n--b--")
	c.Request().SetBody(body)
	c.Request().Header.SetContentLength(len(body))

	val := reflect.ValueOf(&test).Elem()

	err := deserialize(c.RequestContext(), val, nil)
	assert.NoError(suite.T(), err)
}
