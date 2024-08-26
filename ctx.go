package lite

import (
	"context"
	"mime/multipart"
	"net/http"
	"reflect"

	"github.com/go-lite/lite/mime"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type Context[Request any] interface {
	Context() context.Context
	Requests() (Request, error)
	Accepts(offers ...string) string
	AcceptsCharsets(offers ...string) string
	AcceptsEncodings(offers ...string) string
	AcceptsLanguages(offers ...string) string
	App() *App
	Append(field string, values ...string)
	Attachment(filename ...string)
	BaseURL() string
	BodyRaw() []byte
	ClearCookie(key ...string)
	RequestContext() *fasthttp.RequestCtx
	SetUserContext(ctx context.Context)
	Cookie(cookie *fiber.Cookie)
	Cookies(key string, defaultValue ...string) string
	Download(file string, filename ...string) error
	Request() *fasthttp.Request
	Response() *fasthttp.Response
	Get(key string) string
	Format(body interface{}) error
	Hostname() string
	Port() string
	IP() string
	IPs() []string
	Is(extension string) bool
	Links(link ...string)
	Method(override ...string) string
	OriginalURL() string
	SaveFile(fileheader *multipart.FileHeader, path string) error
	Set(key string, val string)
	Status(status int) Context[Request]
	// SetContentType sets the Content-Type response header with the given type and charset.
	SetContentType(extension mime.Mime, charset ...string) Context[Request]
}

var (
	_ Context[string] = &ContextWithRequest[string]{}
	_ Context[any]    = &ContextNoRequest{}
)

type ContextNoRequest struct {
	ctx  *fiber.Ctx
	app  *App
	path string
}

type ContextWithRequest[Request any] struct {
	ContextNoRequest
}

// Context returns the context.
func (c *ContextNoRequest) Context() context.Context {
	return c.ctx.UserContext()
}

func (c *ContextNoRequest) Requests() (any, error) {
	return nil, nil
}

func (c *ContextWithRequest[Request]) Requests() (Request, error) {
	var req Request

	typeOfReq := reflect.TypeOf(&req).Elem()

	reqContext := c.RequestContext()

	params := extractParams(c.path, string(reqContext.Path()))

	switch typeOfReq.Kind() {
	case reflect.Struct:
		err := deserializeRequests(reqContext, &req, params)
		if err != nil {
			return req, err
		}
	case reflect.String:
		err := deserializeBody(reqContext, reflect.ValueOf(&req).Elem())
		if err != nil {
			return req, err
		}
	case reflect.Array, reflect.Slice:
		if typeOfReq.Elem().Kind() == reflect.Uint8 {
			err := deserializeBody(reqContext, reflect.ValueOf(&req).Elem())
			if err != nil {
				return req, err
			}
		} else {
			return req, BadRequestError{
				Context:     "/api/contexts/RequestBodyError",
				Type:        "RequestBodyError",
				Status:      http.StatusBadRequest,
				Title:       "A request body is required",
				Description: "Unsupported slice type",
			}
		}
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		fallthrough
	default:
		return req, BadRequestError{
			Context:     "/api/contexts/RequestBodyError",
			Type:        "RequestBodyError",
			Status:      http.StatusBadRequest,
			Title:       "A request body is required",
			Description: "Unsupported type",
		}
	}

	if typeOfReq.Kind() == reflect.Struct {
		err := c.app.validate(req)
		if err != nil {
			return req, err
		}
	}

	return req, nil
}

func (c *ContextNoRequest) Accepts(offers ...string) string {
	return c.ctx.Accepts(offers...)
}

func (c *ContextNoRequest) AcceptsCharsets(offers ...string) string {
	return c.ctx.AcceptsCharsets(offers...)
}

func (c *ContextNoRequest) AcceptsEncodings(offers ...string) string {
	return c.ctx.AcceptsEncodings(offers...)
}

func (c *ContextNoRequest) AcceptsLanguages(offers ...string) string {
	return c.ctx.AcceptsLanguages(offers...)
}

func (c *ContextNoRequest) App() *App {
	return c.app
}

func (c *ContextNoRequest) Append(field string, values ...string) {
	c.ctx.Append(field, values...)
}

// Attachment adds an attachment to the response.
func (c *ContextNoRequest) Attachment(filename ...string) {
	c.ctx.Attachment(filename...)
}

// BaseURL returns the base URL.
func (c *ContextNoRequest) BaseURL() string {
	return c.ctx.BaseURL()
}

// BodyRaw returns the raw body.
func (c *ContextNoRequest) BodyRaw() []byte {
	return c.ctx.Request().Body()
}

// ClearCookie clears the cookie.
func (c *ContextNoRequest) ClearCookie(key ...string) {
	c.ctx.ClearCookie(key...)
}

// RequestContext returns the request context.
func (c *ContextNoRequest) RequestContext() *fasthttp.RequestCtx {
	return c.ctx.Context()
}

// SetUserContext sets the user context.
func (c *ContextNoRequest) SetUserContext(ctx context.Context) {
	c.ctx.SetUserContext(ctx)
}

// Cookie sets the cookie.
func (c *ContextNoRequest) Cookie(cookie *fiber.Cookie) {
	c.ctx.Cookie(cookie)
}

// Cookies returns the cookie value.
func (c *ContextNoRequest) Cookies(key string, defaultValue ...string) string {
	return c.ctx.Cookies(key, defaultValue...)
}

// Download downloads the file.
func (c *ContextNoRequest) Download(file string, filename ...string) error {
	return c.ctx.Download(file, filename...)
}

// Request returns the request.
func (c *ContextNoRequest) Request() *fasthttp.Request {
	return c.ctx.Request()
}

// Response returns the response.
func (c *ContextNoRequest) Response() *fasthttp.Response {
	return c.ctx.Response()
}

func (c *ContextNoRequest) Get(key string) string {
	return c.ctx.Get(key)
}

// Format formats the response body.
func (c *ContextNoRequest) Format(body interface{}) error {
	return c.ctx.Format(body)
}

// Hostname returns the hostname on which the request is received.
func (c *ContextNoRequest) Hostname() string {
	return c.ctx.Hostname()
}

// Port returns the port on which the request is received.
func (c *ContextNoRequest) Port() string {
	return c.ctx.Port()
}

// IP returns the client IP.
func (c *ContextNoRequest) IP() string {
	return c.ctx.IP()
}

// IPs returns the client IPs.
func (c *ContextNoRequest) IPs() []string {
	return c.ctx.IPs()
}

// Is returns true if the request has the specified extension.
func (c *ContextNoRequest) Is(extension string) bool {
	return c.ctx.Is(extension)
}

// Links adds the specified link to the response.
func (c *ContextNoRequest) Links(link ...string) {
	c.ctx.Links(link...)
}

// Method returns the HTTP method used for the request.
func (c *ContextNoRequest) Method(override ...string) string {
	return c.ctx.Method(override...)
}

// OriginalURL returns the original URL.
func (c *ContextNoRequest) OriginalURL() string {
	return c.ctx.OriginalURL()
}

// SaveFile saves the file to the specified path.
func (c *ContextNoRequest) SaveFile(file *multipart.FileHeader, path string) error {
	return c.ctx.SaveFile(file, path)
}

// Set sets the response's HTTP header field to the specified key, value.
func (c *ContextNoRequest) Set(key string, val string) {
	c.ctx.Set(key, val)
}

// Status sets the HTTP status code.
func (c *ContextNoRequest) Status(status int) Context[any] {
	c.ctx = c.ctx.Status(status)

	return c
}

// Status sets the HTTP status code.
func (c *ContextWithRequest[Request]) Status(status int) Context[Request] {
	c.ctx = c.ctx.Status(status)

	return c
}

// SetContentType sets the Content-Type response header with the given type and charset.
func (c *ContextNoRequest) SetContentType(extension mime.Mime, charset ...string) Context[any] {
	c.ctx = c.ctx.Type(extension, charset...)

	return c
}

// SetContentType sets the Content-Type response header with the given type and charset.
func (c *ContextWithRequest[Request]) SetContentType(extension mime.Mime, charset ...string) Context[Request] {
	c.ctx = c.ctx.Type(extension, charset...)

	return c
}
