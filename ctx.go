package lite

import (
	"context"
	"crypto/tls"
	"errors"
	"mime/multipart"
	"reflect"

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
	Format(body interface{}) error
	Fresh() bool
	Hostname() string
	Port() string
	IP() string
	IPs() []string
	Is(extension string) bool
	Links(link ...string)
	Method(override ...string) string
	MultipartForm() (*multipart.Form, error)
	ClientHelloInfo() *tls.ClientHelloInfo
	Next() error
	OriginalURL() string
	Protocol() string
	SaveFile(fileheader *multipart.FileHeader, path string) error
	Set(key string, val string)
	Status(status int) Context[Request]
	// SetContentType sets the Content-Type response header with the given type and charset.
	SetContentType(extension string, charset ...string) Context[Request]
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

func (c *ContextNoRequest) Context() context.Context {
	return c.ctx.UserContext()
}

func (c *ContextNoRequest) Requests() (any, error) {
	var req any

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
			err := deserialize(reqContext, reflect.ValueOf(&req).Elem(), params)
			if err != nil {
				return req, err
			}
		} else {
			return req, errors.New("unsupported slice type")
		}
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		fallthrough
	default:
		return req, errors.New("unsupported type")
	}

	return req, nil
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
			err := deserialize(reqContext, reflect.ValueOf(&req).Elem(), params)
			if err != nil {
				return req, err
			}
		} else {
			return req, errors.New("unsupported slice type")
		}
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		fallthrough
	default:
		return req, errors.New("unsupported type")
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

func (c *ContextNoRequest) Attachment(filename ...string) {
	c.ctx.Attachment(filename...)
}

func (c *ContextNoRequest) BaseURL() string {
	return c.ctx.BaseURL()
}

func (c *ContextNoRequest) BodyRaw() []byte {
	return c.ctx.Request().Body()
}

func (c *ContextNoRequest) ClearCookie(key ...string) {
	c.ctx.ClearCookie(key...)
}

func (c *ContextNoRequest) RequestContext() *fasthttp.RequestCtx {
	return c.ctx.Context()
}

func (c *ContextNoRequest) SetUserContext(ctx context.Context) {
	c.ctx.SetUserContext(ctx)
}

func (c *ContextNoRequest) Cookie(cookie *fiber.Cookie) {
	c.ctx.Cookie(cookie)
}

func (c *ContextNoRequest) Cookies(key string, defaultValue ...string) string {
	return c.ctx.Cookies(key, defaultValue...)
}

func (c *ContextNoRequest) Download(file string, filename ...string) error {
	return c.ctx.Download(file, filename...)
}

func (c *ContextNoRequest) Request() *fasthttp.Request {
	return c.ctx.Request()
}

func (c *ContextNoRequest) Response() *fasthttp.Response {
	return c.ctx.Response()
}

func (c *ContextNoRequest) Format(body interface{}) error {
	return c.ctx.Format(body)
}

func (c *ContextNoRequest) Fresh() bool {
	return c.ctx.Fresh()
}

func (c *ContextNoRequest) Hostname() string {
	return c.ctx.Hostname()
}

func (c *ContextNoRequest) Port() string {
	return c.ctx.Port()
}

func (c *ContextNoRequest) IP() string {
	return c.ctx.IP()
}

func (c *ContextNoRequest) IPs() []string {
	return c.ctx.IPs()
}

func (c *ContextNoRequest) Is(extension string) bool {
	return c.ctx.Is(extension)
}

func (c *ContextNoRequest) Links(link ...string) {
	c.ctx.Links(link...)
}

func (c *ContextNoRequest) Method(override ...string) string {
	return c.ctx.Method(override...)
}

func (c *ContextNoRequest) MultipartForm() (*multipart.Form, error) {
	return c.ctx.MultipartForm()
}

func (c *ContextNoRequest) ClientHelloInfo() *tls.ClientHelloInfo {
	return c.ctx.ClientHelloInfo()
}

func (c *ContextNoRequest) Next() error {
	return c.ctx.Next()
}

func (c *ContextNoRequest) OriginalURL() string {
	return c.ctx.OriginalURL()
}

func (c *ContextNoRequest) Protocol() string {
	return c.ctx.Protocol()
}

func (c *ContextNoRequest) SaveFile(file *multipart.FileHeader, path string) error {
	return c.ctx.SaveFile(file, path)
}

// Set sets the response's HTTP header field to the specified key, value.
func (c *ContextNoRequest) Set(key string, val string) {
	c.ctx.Set(key, val)
}

func (c *ContextNoRequest) Status(status int) Context[any] {
	c.ctx = c.ctx.Status(status)

	return c
}

func (c *ContextWithRequest[Request]) Status(status int) Context[Request] {
	c.ctx = c.ctx.Status(status)

	return c
}

func (c *ContextNoRequest) SetContentType(extension string, charset ...string) Context[any] {
	c.ctx = c.ctx.Type(extension, charset...)

	return c
}

func (c *ContextWithRequest[Request]) SetContentType(extension string, charset ...string) Context[Request] {
	c.ctx = c.ctx.Type(extension, charset...)

	return c
}
