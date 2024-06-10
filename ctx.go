package lite

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/disco07/lite-fiber/lite"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"reflect"
)

type Context[Request any] interface {
	Context() context.Context
	Requests() (Request, error)
	Accepts(offers ...string) string
	AcceptsCharsets(offers ...string) string
	AcceptsEncodings(offers ...string) string
	AcceptsLanguages(offers ...string) string
	App() *lite.App
	Append(field string, values ...string)
	Attachment(filename ...string)
	BaseURL() string
	BodyRaw() []byte
	tryDecodeBodyInOrder(originalBody *[]byte, encodings []string) ([]byte, uint8, error)
	Body() []byte
	BodyParser(out interface{}) error
	ClearCookie(key ...string)
	RequestContext() *fasthttp.RequestCtx
	SetUserContext(ctx context.Context)
	Cookie(cookie *fiber.Cookie)
	Cookies(key string, defaultValue ...string) string
	CookieParser(out interface{}) error
	Download(file string, filename ...string) error
	Request() *fasthttp.Request
	Response() *fasthttp.Response
	Format(body interface{}) error
	FormFile(key string) (*multipart.FileHeader, error)
	FormValue(key string, defaultValue ...string) string
	Fresh() bool
	Get(key string, defaultValue ...string) string
	GetRespHeader(key string, defaultValue ...string) string
	GetReqHeaders() map[string][]string
	GetRespHeaders() map[string][]string
	Hostname() string
	Port() string
	IP() string
	IPs() []string
	Is(extension string) bool
	JSON(data interface{}, ctype ...string) error
	JSONP(data interface{}, callback ...string) error
	XML(data interface{}) error
	Links(link ...string)
	Locals(key interface{}, value ...interface{}) interface{}
	Location(path string)
	Method(override ...string) string
	MultipartForm() (*multipart.Form, error)
	ClientHelloInfo() *tls.ClientHelloInfo
	Next() error
	RestartRouting() error
	OriginalURL() string
	Params(key string, defaultValue ...string) string
	AllParams() map[string]string
	ParamsParser(out interface{}) error
	ParamsInt(key string, defaultValue ...int) (int, error)
	Path(override ...string) string
	Protocol() string
	Query(key string, defaultValue ...string) string
	Queries() map[string]string
	QueryInt(key string, defaultValue ...int) int
	QueryBool(key string, defaultValue ...bool) bool
	QueryFloat(key string, defaultValue ...float64) float64
	QueryParser(out interface{}) error
	ReqHeaderParser(out interface{}) error
	parseToStruct(aliasTag string, out interface{}, data map[string][]string) error
	Range(size int) (fiber.Range, error)
	Redirect(location string, status ...int) error
	Bind(vars fiber.Map) error
	GetRouteURL(routeName string, params fiber.Map) (string, error)
	RedirectToRoute(routeName string, params fiber.Map, status ...int) error
	RedirectBack(fallback string, status ...int) error
	Render(name string, bind interface{}, layouts ...string) error
	Route() *fiber.Route
	SaveFile(fileheader *multipart.FileHeader, path string) error
	SaveFileToStorage(fileheader *multipart.FileHeader, path string, storage fiber.Storage) error
	Secure() bool
	Send(body []byte) error
	SendFile(file string, compress ...bool) error
	SendStatus(status int) error
	SendString(body string) error
	SendStream(stream io.Reader, size ...int) error
	Set(key string, val string)
	Subdomains(offset ...int) []string
	Stale() bool
	Status(status int) *lite.Context[Request]
	String() string
	Type(extension string, charset ...string) *lite.Context[Request]
	Vary(fields ...string)
	Write(p []byte) (int, error)
	Writef(f string, a ...interface{}) (int, error)
	WriteString(s string) (int, error)
	XHR() bool
	IsProxyTrusted() bool
	IsFromLocal() bool
}

type ContextWithRequest[Request any] struct {
	ctx  *fiber.Ctx
	app  *lite.App
	path string
}

func (c *ContextWithRequest[Request]) Context() context.Context {
	return c.ctx.UserContext()
}

func (c *ContextWithRequest[Request]) Requests() (Request, error) {
	var req Request

	typeOfReq := reflect.TypeOf(&req).Elem()

	reqContext := c.RequestContext()

	params := extractParams(c.path, string(reqContext.Path()))

	switch typeOfReq.Kind() {
	case reflect.Struct:
		err := deserializeParams(reqContext, &req, params)
		if err != nil {
			return req, err
		}

		err = deserializeBody(reqContext, &req)
		if err != nil {
			return req, err
		}
	default:
		return req, errors.New("unsupported type")
	}

	return req, nil
}

func (c *ContextWithRequest[Request]) Accepts(offers ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) AcceptsCharsets(offers ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) AcceptsEncodings(offers ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) AcceptsLanguages(offers ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) App() *lite.App {
	return c.app
}

func (c *ContextWithRequest[Request]) Append(field string, values ...string) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Attachment(filename ...string) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) BaseURL() string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) BodyRaw() []byte {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) tryDecodeBodyInOrder(originalBody *[]byte, encodings []string) ([]byte, uint8, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Body() []byte {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) BodyParser(out interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) ClearCookie(key ...string) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) RequestContext() *fasthttp.RequestCtx {
	return c.ctx.Context()
}

func (c *ContextWithRequest[Request]) SetUserContext(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Cookie(cookie *fiber.Cookie) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Cookies(key string, defaultValue ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) CookieParser(out interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Download(file string, filename ...string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Request() *fasthttp.Request {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Response() *fasthttp.Response {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Format(body interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) FormFile(key string) (*multipart.FileHeader, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) FormValue(key string, defaultValue ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Fresh() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Get(key string, defaultValue ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) GetRespHeader(key string, defaultValue ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) GetReqHeaders() map[string][]string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) GetRespHeaders() map[string][]string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Hostname() string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Port() string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) IP() string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) IPs() []string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Is(extension string) bool {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) JSON(data interface{}, ctype ...string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) JSONP(data interface{}, callback ...string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) XML(data interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Links(link ...string) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Locals(key interface{}, value ...interface{}) interface{} {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Location(path string) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Method(override ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) MultipartForm() (*multipart.Form, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) ClientHelloInfo() *tls.ClientHelloInfo {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Next() error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) RestartRouting() error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) OriginalURL() string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Params(key string, defaultValue ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) AllParams() map[string]string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) ParamsParser(out interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) ParamsInt(key string, defaultValue ...int) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Path(override ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Protocol() string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Query(key string, defaultValue ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Queries() map[string]string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) QueryInt(key string, defaultValue ...int) int {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) QueryBool(key string, defaultValue ...bool) bool {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) QueryFloat(key string, defaultValue ...float64) float64 {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) QueryParser(out interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) ReqHeaderParser(out interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) parseToStruct(aliasTag string, out interface{}, data map[string][]string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Range(size int) (fiber.Range, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Redirect(location string, status ...int) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Bind(vars fiber.Map) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) GetRouteURL(routeName string, params fiber.Map) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) RedirectToRoute(routeName string, params fiber.Map, status ...int) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) RedirectBack(fallback string, status ...int) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Render(name string, bind interface{}, layouts ...string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) renderExtensions(bind interface{}) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Route() *fiber.Route {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) SaveFile(fileheader *multipart.FileHeader, path string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) SaveFileToStorage(fileheader *multipart.FileHeader, path string, storage fiber.Storage) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Secure() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Send(body []byte) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) SendFile(file string, compress ...bool) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) SendStatus(status int) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) SendString(body string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) SendStream(stream io.Reader, size ...int) error {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Set(key string, val string) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Subdomains(offset ...int) []string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Stale() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Status(status int) *lite.Context[Request] {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) String() string {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Type(extension string, charset ...string) *lite.Context[Request] {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Vary(fields ...string) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Write(p []byte) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) Writef(f string, a ...interface{}) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) WriteString(s string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) XHR() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) IsProxyTrusted() bool {
	//TODO implement me
	panic("implement me")
}

func (c *ContextWithRequest[Request]) IsFromLocal() bool {
	//TODO implement me
	panic("implement me")
}
