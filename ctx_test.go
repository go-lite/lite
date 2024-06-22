package lite

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
	"strings"
	"testing"
	"time"
)

type CtxTestSuite struct {
	suite.Suite
}

func (suite *CtxTestSuite) SetupTest() {
}

func TestCtxTestSuite(t *testing.T) {
	suite.Run(t, new(CtxTestSuite))
}

func newContext[Request any](ctx *fiber.Ctx, app *App, path string) Context[Request] {
	c := ContextNoRequest{
		ctx:  ctx,
		app:  app,
		path: path,
	}

	return &ContextWithRequest[Request]{
		ContextNoRequest: c,
	}
}

func (suite *CtxTestSuite) TestContextWithRequest_Body_TextPlain() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "text/plain")
	ctx.Request().SetBodyString("Hello World")

	c := newContext[string](ctx, app, "/foo")
	req, err := c.Requests()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Hello World", req)
}

type request struct {
	ID     uint64   `lite:"path=id"`
	Header bool     `lite:"header=X-Real-Ip"`
	Q      bool     `lite:"query=q"`
	Body   bodyTest `lite:"req=body"`
}

type requestMultiParams struct {
	P pathParams
}

type pathParams struct {
	ID   uint64 `lite:"path=id"`
	Name string `lite:"path=name"`
}

type bodyTest struct {
	A string `json:"A" yaml:"A" xml:"A"`
	B int    `json:"B" yaml:"B" xml:"B"`
	C bool   `json:"C" yaml:"C" xml:"C"`
}

func (suite *CtxTestSuite) TestContextWithRequest_ApplicationJSON_Requests() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":true}`)

	c := newContext[request](ctx, app, "/foo")
	req, err := c.Requests()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), bodyTest{A: "a", B: 1, C: true}, req.Body)
}

func (suite *CtxTestSuite) TestContextWithRequest_ApplicationXML_Requests() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/xml")
	ctx.Request().SetBodyString(`<request><A>a</A><B>1</B><C>true</C></request>`)

	c := newContext[request](ctx, app, "/foo")
	req, err := c.Requests()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), bodyTest{A: "a", B: 1, C: true}, req.Body)
}

func (suite *CtxTestSuite) TestContextWithRequest_ApplicationJSON_Requests_Error() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":"1"}`)

	c := newContext[request](ctx, app, "/foo")
	_, err := c.Requests()
	assert.Error(suite.T(), err)
}

func (suite *CtxTestSuite) TestContextWithRequest_Path() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/foo/42")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":true}`)

	c := newContext[request](ctx, app, "/foo/:id")
	req, err := c.Requests()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint64(42), req.ID)
}

func (suite *CtxTestSuite) TestContextWithRequest_MultiPath() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/foo/42/john")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":true}`)

	c := newContext[requestMultiParams](ctx, app, "/foo/:id/:name")
	req, err := c.Requests()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint64(42), req.P.ID)
	assert.Equal(suite.T(), "john", req.P.Name)
}

func (suite *CtxTestSuite) TestContextWithRequest_MultiPath_Error() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/foo/doe/john")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":true}`)

	c := newContext[requestMultiParams](ctx, app, "/foo/:id/:name")
	_, err := c.Requests()
	assert.Error(suite.T(), err)
}

func (suite *CtxTestSuite) TestContextWithRequest_Path_Error() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/foo/abc")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":true}`)

	c := newContext[request](ctx, app, "/foo/:id")
	_, err := c.Requests()
	assert.Error(suite.T(), err)
}

func (suite *CtxTestSuite) TestContextWithRequest_Missing_Path_Error() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/foo/abc")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":true}`)

	c := newContext[request](ctx, app, "/foo/:name")
	_, err := c.Requests()
	assert.Error(suite.T(), err)
}

func (suite *CtxTestSuite) TestContextWithRequest_Query() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/foo/42?q=true")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":true}`)

	c := newContext[request](ctx, app, "/foo/:id")
	req, err := c.Requests()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, req.Q)
}

func (suite *CtxTestSuite) TestContextWithRequest_Query_Error() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/foo/42?q=hello")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":true}`)

	c := newContext[request](ctx, app, "/foo/:id")
	_, err := c.Requests()
	assert.Error(suite.T(), err)
}

func (suite *CtxTestSuite) TestContextWithRequest_Header() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().Header.Set("X-Real-Ip", "true")
	ctx.Request().SetRequestURI("/foo/42?q=true")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":true}`)

	c := newContext[request](ctx, app, "/foo/:id")
	req, err := c.Requests()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, req.Header)
}

func (suite *CtxTestSuite) TestContextWithRequest_Header_Error() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")
	ctx.Request().Header.Set("X-Real-Ip", "test")
	ctx.Request().SetRequestURI("/foo/42?q=true")
	ctx.Request().SetBodyString(`{"A":"a","B":1,"C":true}`)

	c := newContext[request](ctx, app, "/foo/:id")
	_, err := c.Requests()
	assert.Error(suite.T(), err)
}

func (suite *CtxTestSuite) TestContextWithRequest_Error() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[chan int](ctx, app, "/foo/:id")
	_, err := c.Requests()
	assert.Error(suite.T(), err)
}

func (suite *CtxTestSuite) TestContextWithRequest_Accepts() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Accept", "application/json")

	c := newContext[request](ctx, app, "/foo")
	assert.True(suite.T(), "application/json" == c.Accepts("application/json"))
	assert.False(suite.T(), "application/json" == c.Accepts("application/xml"))
}

func (suite *CtxTestSuite) TestContextWithRequest_AcceptsCharsets() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Accept-Charset", "utf-8")

	c := newContext[request](ctx, app, "/foo")
	assert.True(suite.T(), "utf-8" == c.AcceptsCharsets("utf-8"))
	assert.False(suite.T(), "utf-8" == c.AcceptsCharsets("utf-16"))
}

func (suite *CtxTestSuite) TestContextWithRequest_AcceptsEncodings() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Accept-Encoding", "gzip")

	c := newContext[request](ctx, app, "/foo")
	assert.True(suite.T(), "gzip" == c.AcceptsEncodings("gzip"))
	assert.False(suite.T(), "gzip" == c.AcceptsEncodings("deflate"))
}

func (suite *CtxTestSuite) TestContextWithRequest_AcceptsLanguages() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Accept-Language", "en")

	c := newContext[request](ctx, app, "/foo")
	assert.True(suite.T(), "en" == c.AcceptsLanguages("en"))
	assert.False(suite.T(), "en" == c.AcceptsLanguages("fr"))
}

func (suite *CtxTestSuite) TestContextWithRequest_Is() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Content-Type", "application/json")

	c := newContext[request](ctx, app, "/foo")
	assert.True(suite.T(), c.Is("json"))
	assert.False(suite.T(), c.Is("xml"))
}

func (suite *CtxTestSuite) TestContextWithRequest_App() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[request](ctx, app, "/foo")
	assert.Equal(suite.T(), app, c.App())
}

func (suite *CtxTestSuite) TestContextWithRequest_OriginalURL() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().SetRequestURI("/foo")

	c := newContext[request](ctx, app, "/foo")
	assert.Equal(suite.T(), "/foo", c.OriginalURL())
}

func (suite *CtxTestSuite) TestContextWithRequest_Append() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.Set("Authorization", "test")

	c := newContext[request](ctx, app, "/foo")
	c.Append("Authorization", "test")
	assert.Equal(suite.T(), "test", string(ctx.Request().Header.Peek("Authorization")))
}

func (suite *CtxTestSuite) TestContextWithRequest_Attachment() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[request](ctx, app, "/foo")
	c.Attachment("test")
	assert.Equal(suite.T(), "attachment; filename=\"test\"", string(ctx.Response().Header.Peek("Content-Disposition")))
}

func (suite *CtxTestSuite) TestContextWithRequest_BaseURL() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().SetRequestURI("http://www.test.com/")

	c := newContext[request](ctx, app, "/foo")
	assert.Equal(suite.T(), "http://www.test.com", c.BaseURL())
}

func (suite *CtxTestSuite) TestContextWithRequest_BodyRaw() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().SetBodyString("test")

	c := newContext[request](ctx, app, "/foo")
	assert.Equal(suite.T(), "test", string(c.BodyRaw()))
}

func (suite *CtxTestSuite) TestContextWithRequest_ClearCookie() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.SetCookie("test", "test")

	c := newContext[request](ctx, app, "/foo")
	c.ClearCookie("test")
	assert.Equal(suite.T(), "test=; expires=Tue, 10 Nov 2009 23:00:00 GMT", string(ctx.Response().Header.Peek("Set-Cookie")))
}

func (suite *CtxTestSuite) TestContextWithRequest_RequestContext() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[request](ctx, app, "/foo")
	assert.Equal(suite.T(), ctx.Context(), c.RequestContext())
}

func (suite *CtxTestSuite) TestContextWithRequest_SetUserContext() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[request](ctx, app, "/foo")
	c.SetUserContext(context.Background())
	assert.Equal(suite.T(), context.Background(), c.Context())
}

func (suite *CtxTestSuite) TestContextWithRequest_Cookie() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.SetCookie("test", "test")

	c := newContext[request](ctx, app, "/foo")

	expire := time.Now().Add(24 * time.Hour)
	var dst []byte
	dst = expire.In(time.UTC).AppendFormat(dst, time.RFC1123)
	httpdate := strings.ReplaceAll(string(dst), "UTC", "GMT")

	c.Cookie(&fiber.Cookie{
		Name:    "username",
		Value:   "john",
		Expires: expire,
	})

	assert.Equal(suite.T(), "username=john; expires="+httpdate+"; path=/; SameSite=Lax", string(ctx.Response().Header.Peek("Set-Cookie")))
}

func (suite *CtxTestSuite) TestContextWithRequest_Cookies() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	ctx.Request().Header.SetCookie("test", "test")

	c := newContext[request](ctx, app, "/foo")

	assert.Equal(suite.T(), "test", c.Cookies("test"))
}

func (suite *CtxTestSuite) TestContextWithRequestDownload() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[request](ctx, app, "/foo")

	err := c.Download("ctx.go")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "attachment; filename=\"ctx.go\"", string(ctx.Response().Header.Peek("Content-Disposition")))
}

func (suite *CtxTestSuite) TestContextWithRequest_Request() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[request](ctx, app, "/foo")

	assert.Equal(suite.T(), ctx.Request(), c.Request())
}

func (suite *CtxTestSuite) TestContextWithRequest_Response() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[request](ctx, app, "/foo")

	assert.Equal(suite.T(), ctx.Response(), c.Response())
}

func (suite *CtxTestSuite) TestContextWithRequest_Format() {
	app := New()

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	c := newContext[request](ctx, app, "/foo")

	err := c.Format("json")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "text/html", string(ctx.Response().Header.Peek("Content-Type")))
}
