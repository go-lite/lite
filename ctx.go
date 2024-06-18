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
	FormFile(key string) (*multipart.FileHeader, error)
	FormValue(key string, defaultValue ...string) string
	Fresh() bool
	Get(key string, defaultValue ...string) string
	GetReqHeaders() map[string][]string
	GetRespHeaders() map[string][]string
	Hostname() string
	Port() string
	IP() string
	IPs() []string
	Is(extension string) bool
	Links(link ...string)
	Locals(key interface{}, value ...interface{}) interface{}
	Location(path string)
	Method(override ...string) string
	MultipartForm() (*multipart.Form, error)
	ClientHelloInfo() *tls.ClientHelloInfo
	Next() error
	RestartRouting() error
	OriginalURL() string
	AllParams() map[string]string
	Path(override ...string) string
	Protocol() string
	Queries() map[string]string
	Range(size int) (fiber.Range, error)
	Redirect(location string, status ...int) error
	Bind(vars fiber.Map) error
	GetRouteURL(routeName string, params fiber.Map) (string, error)
	RedirectToRoute(routeName string, params fiber.Map, status ...int) error
	RedirectBack(fallback string, status ...int) error
	Render(name string, bind interface{}, layouts ...string) error
	SaveFile(fileheader *multipart.FileHeader, path string) error
	SaveFileToStorage(fileheader *multipart.FileHeader, path string, storage fiber.Storage) error
	Secure() bool
	Set(key string, val string)
	Subdomains(offset ...int) []string
	Stale() bool
	Status(status int) Context[Request]
	String() string
	// Type sets the Content-Type response header with the given type and charset.
	Type(extension string, charset ...string) Context[Request]
	Vary(fields ...string)
	XHR() bool
	IsProxyTrusted() bool
	IsFromLocal() bool
}

var (
	_ Context[string] = &ContextWithRequest[string]{}
	_ Context[any]    = &ContextNoRequest{}
)

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

func (c *ContextNoRequest) FormFile(key string) (*multipart.FileHeader, error) {
	return c.ctx.FormFile(key)
}

func (c *ContextNoRequest) FormValue(key string, defaultValue ...string) string {
	return c.ctx.FormValue(key, defaultValue...)
}

func (c *ContextNoRequest) Fresh() bool {
	return c.ctx.Fresh()
}

func (c *ContextNoRequest) Get(key string, defaultValue ...string) string {
	return c.ctx.Get(key, defaultValue...)
}

func (c *ContextNoRequest) GetReqHeaders() map[string][]string {
	return c.ctx.GetReqHeaders()
}

func (c *ContextNoRequest) GetRespHeaders() map[string][]string {
	return c.ctx.GetRespHeaders()
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

func (c *ContextNoRequest) Locals(key interface{}, value ...interface{}) interface{} {
	return c.ctx.Locals(key, value...)
}

func (c *ContextNoRequest) Location(path string) {
	c.ctx.Location(path)
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

func (c *ContextNoRequest) RestartRouting() error {
	return c.ctx.RestartRouting()
}

func (c *ContextNoRequest) OriginalURL() string {
	return c.ctx.OriginalURL()
}

func (c *ContextNoRequest) Params(key string, defaultValue ...string) string {
	return c.ctx.Params(key, defaultValue...)
}

func (c *ContextNoRequest) AllParams() map[string]string {
	return c.ctx.AllParams()
}

func (c *ContextNoRequest) Path(override ...string) string {
	return c.ctx.Path(override...)
}

func (c *ContextNoRequest) Protocol() string {
	return c.ctx.Protocol()
}

func (c *ContextNoRequest) Queries() map[string]string {
	return c.ctx.Queries()
}

func (c *ContextNoRequest) QueryInt(key string, defaultValue ...int) int {
	return c.ctx.QueryInt(key, defaultValue...)
}

func (c *ContextNoRequest) Range(size int) (fiber.Range, error) {
	return c.ctx.Range(size)
}

func (c *ContextNoRequest) Redirect(location string, status ...int) error {
	return c.ctx.Redirect(location, status...)
}

func (c *ContextNoRequest) Bind(vars fiber.Map) error {
	return c.ctx.Bind(vars)
}

func (c *ContextNoRequest) GetRouteURL(routeName string, params fiber.Map) (string, error) {
	return c.ctx.GetRouteURL(routeName, params)
}

func (c *ContextNoRequest) RedirectToRoute(routeName string, params fiber.Map, status ...int) error {
	return c.ctx.RedirectToRoute(routeName, params, status...)
}

func (c *ContextNoRequest) RedirectBack(fallback string, status ...int) error {
	return c.ctx.RedirectBack(fallback, status...)
}

func (c *ContextNoRequest) Render(name string, bind interface{}, layouts ...string) error {
	return c.ctx.Render(name, bind, layouts...)
}

func (c *ContextNoRequest) SaveFile(file *multipart.FileHeader, path string) error {
	return c.ctx.SaveFile(file, path)
}

func (c *ContextNoRequest) SaveFileToStorage(fileheader *multipart.FileHeader, path string, storage fiber.Storage) error {
	return c.ctx.SaveFileToStorage(fileheader, path, storage)
}

func (c *ContextNoRequest) Secure() bool {
	return c.ctx.Secure()
}

// Set sets the response's HTTP header field to the specified key, value.
func (c *ContextNoRequest) Set(key string, val string) {
	c.ctx.Set(key, val)
}

func (c *ContextNoRequest) Subdomains(offset ...int) []string {
	return c.ctx.Subdomains(offset...)
}

func (c *ContextNoRequest) Stale() bool {
	return c.ctx.Stale()
}

func (c *ContextNoRequest) Status(status int) Context[any] {
	c.ctx = c.ctx.Status(status)

	return c
}

func (c *ContextWithRequest[Request]) Status(status int) Context[Request] {
	c.ctx = c.ctx.Status(status)

	return c
}

func (c *ContextNoRequest) String() string {
	return c.ctx.String()
}

func (c *ContextNoRequest) Type(extension string, charset ...string) Context[any] {
	c.ctx = c.ctx.Type(extension, charset...)

	return c
}

func (c *ContextWithRequest[Request]) Type(extension string, charset ...string) Context[Request] {
	c.ctx = c.ctx.Type(extension, charset...)

	return c
}

func (c *ContextNoRequest) Vary(fields ...string) {
	c.ctx.Vary(fields...)
}

func (c *ContextNoRequest) XHR() bool {
	return c.ctx.XHR()
}

func (c *ContextNoRequest) IsProxyTrusted() bool {
	return c.ctx.IsProxyTrusted()
}

func (c *ContextNoRequest) IsFromLocal() bool {
	return c.ctx.IsFromLocal()
}
