package lite

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	liteErrors "github.com/go-lite/lite/errors"
	"github.com/gofiber/fiber/v2"
)

func newLiteContext[Request any, Contexter Context[Request]](ctx ContextNoRequest) Contexter {
	var c Contexter

	switch any(c).(type) {
	case *ContextNoRequest:
		return any(&ctx).(Contexter)
	case *ContextWithRequest[Request]:
		return any(&ContextWithRequest[Request]{
			ContextNoRequest: ctx,
		}).(Contexter)
	default:
		panic("unknown type")
	}
}

func fiberHandler[ResponseBody, Request any, Contexter Context[Request]](
	controller func(c Contexter) (ResponseBody, error),
	path string,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Context().SetContentType("application/json")

		ctx := newLiteContext[Request, Contexter](ContextNoRequest{ctx: c, path: path})

		c.Status(getStatusCode(c.Method()))

		response, err := controller(ctx)
		if err != nil {
			// check if the error is a HTTPError and if so, return the error code
			var httpError liteErrors.HTTPError
			if errors.As(err, &httpError) {
				c.Status(httpError.StatusCode())

				return c.JSON(httpError)
			}

			c.Status(http.StatusInternalServerError)

			return c.JSON(liteErrors.DefaultErrorResponses[http.StatusInternalServerError].SetMessage(err.Error()))
		}

		return serializeResponse(c.Context(), &response)
	}
}

func Group(app *App, path string) *App {
	path = strings.TrimRight(path, "/")

	a := *app
	newApp := &a
	newApp.basePath += path
	newApp.tag = toTitle(strings.TrimLeft(path, "/"))

	return newApp
}

func Use(app *App, args ...any) {
	app.app.Use(args...)
}

func Get[ResponseBody, Request any, Contexter Context[Request]](
	app *App,
	path string,
	controller func(Contexter) (ResponseBody, error),
	middleware ...fiber.Handler,
) Route[ResponseBody, Request] {
	return registerRoute[ResponseBody, Request](
		app,
		Route[ResponseBody, Request]{
			path:        path,
			method:      http.MethodGet,
			contentType: "application/json",
			statusCode:  getStatusCode(http.MethodGet),
		},
		fiberHandler[ResponseBody, Request](controller, path),
		middleware...,
	)
}

func Post[ResponseBody, Request any, Contexter Context[Request]](
	app *App,
	path string,
	controller func(Contexter) (ResponseBody, error),
	middleware ...fiber.Handler,
) Route[ResponseBody, Request] {
	return registerRoute[ResponseBody, Request](
		app,
		Route[ResponseBody, Request]{
			path:        path,
			method:      http.MethodPost,
			contentType: "application/json",
			statusCode:  getStatusCode(http.MethodPost),
		},
		fiberHandler[ResponseBody, Request](controller, path),
		middleware...,
	)
}

func Put[ResponseBody, Request any, Contexter Context[Request]](
	app *App,
	path string,
	controller func(Contexter) (ResponseBody, error),
	middleware ...fiber.Handler,
) Route[ResponseBody, Request] {
	return registerRoute[ResponseBody, Request](
		app,
		Route[ResponseBody, Request]{
			path:        path,
			method:      http.MethodPut,
			contentType: "application/json",
			statusCode:  getStatusCode(http.MethodPut),
		},
		fiberHandler[ResponseBody, Request](controller, path),
		middleware...,
	)
}

func Delete[ResponseBody, Request any, Contexter Context[Request]](
	app *App,
	path string,
	controller func(Contexter) (ResponseBody, error),
	middleware ...fiber.Handler,
) Route[ResponseBody, Request] {
	return registerRoute[ResponseBody, Request](
		app,
		Route[ResponseBody, Request]{
			path:        path,
			method:      http.MethodDelete,
			contentType: "application/json",
			statusCode:  getStatusCode(http.MethodDelete),
		},
		fiberHandler[ResponseBody, Request](controller, path),
		middleware...,
	)
}

func Patch[ResponseBody, Request any, Contexter Context[Request]](
	app *App,
	path string,
	controller func(Contexter) (ResponseBody, error),
	middleware ...fiber.Handler,
) Route[ResponseBody, Request] {
	return registerRoute[ResponseBody, Request](
		app,
		Route[ResponseBody, Request]{
			path:        path,
			method:      http.MethodPatch,
			contentType: "application/json",
			statusCode:  getStatusCode(http.MethodPatch),
		},
		fiberHandler[ResponseBody, Request](controller, path),
		middleware...,
	)
}

func Head[ResponseBody, Request any, Contexter Context[Request]](
	app *App,
	path string,
	controller func(Contexter) (ResponseBody, error),
	middleware ...fiber.Handler,
) Route[ResponseBody, Request] {
	return registerRoute[ResponseBody, Request](
		app,
		Route[ResponseBody, Request]{
			path:        path,
			method:      http.MethodHead,
			contentType: "application/json",
			statusCode:  getStatusCode(http.MethodHead),
		},
		fiberHandler[ResponseBody, Request](controller, path),
		middleware...,
	)
}

func Connect[ResponseBody, Request any, Contexter Context[Request]](
	app *App,
	path string,
	controller func(Contexter) (ResponseBody, error),
	middleware ...fiber.Handler,
) Route[ResponseBody, Request] {
	return registerRoute[ResponseBody, Request](
		app,
		Route[ResponseBody, Request]{
			path:        path,
			method:      http.MethodConnect,
			contentType: "application/json",
			statusCode:  getStatusCode(http.MethodConnect),
		},
		fiberHandler[ResponseBody, Request](controller, path),
		middleware...,
	)
}

func Trace[ResponseBody, Request any, Contexter Context[Request]](
	app *App,
	path string,
	controller func(Contexter) (ResponseBody, error),
	middleware ...fiber.Handler,
) Route[ResponseBody, Request] {
	return registerRoute[ResponseBody, Request](
		app,
		Route[ResponseBody, Request]{
			path:        path,
			method:      http.MethodTrace,
			contentType: "application/json",
			statusCode:  getStatusCode(http.MethodTrace),
		},
		fiberHandler[ResponseBody, Request](controller, path),
		middleware...,
	)
}

func Options[ResponseBody, Request any, Contexter Context[Request]](
	app *App,
	path string,
	controller func(Contexter) (ResponseBody, error),
	middleware ...fiber.Handler,
) Route[ResponseBody, Request] {
	return registerRoute[ResponseBody, Request](
		app,
		Route[ResponseBody, Request]{
			path:        path,
			method:      http.MethodOptions,
			contentType: "application/json",
			statusCode:  getStatusCode(http.MethodOptions),
		},
		fiberHandler[ResponseBody, Request](controller, path),
		middleware...,
	)
}

func registerRoute[ResponseBody, Request any](
	app *App,
	route Route[ResponseBody, Request],
	controller fiber.Handler,
	middleware ...fiber.Handler,
) Route[ResponseBody, Request] {
	fullPath := app.basePath + route.path

	if len(middleware) > 0 {
		app.app.Add(
			route.method,
			fullPath,
			middleware...,
		)
	}

	app.app.Add(
		route.method,
		fullPath,
		controller,
	)

	operation, err := registerOpenAPIOperation[ResponseBody, Request](
		app,
		route.method,
		fullPath,
		route.contentType,
		route.statusCode,
	)
	if err != nil {
		slog.ErrorContext(context.Background(), "failed to register openapi operation", slog.Any("error", err))
		panic(err)
	}

	if app.tag != "" {
		operation.Tags = append(operation.Tags, app.tag)
		operation.Description = setDescription(route.method, app.tag)
	}

	operation.OperationID = route.path

	route.operation = operation

	return route
}

// setDescription set the description of the route based on the method and tag
func setDescription(method string, tag string) (description string) {
	switch method {
	case http.MethodGet:
		description = fmt.Sprintf("Get the %s resource", tag)
	case http.MethodPost:
		description = fmt.Sprintf("Create a new %s resource", tag)
	case http.MethodPut:
		description = fmt.Sprintf("Replace the %s resource", tag)
	case http.MethodPatch:
		description = fmt.Sprintf("Update the %s resource", tag)
	case http.MethodDelete:
		description = fmt.Sprintf("Delete the %s resource", tag)
	case http.MethodHead:
		description = fmt.Sprintf("Get the %s resource header", tag)
	case http.MethodOptions:
		description = fmt.Sprintf("Get the %s resource options", tag)
	case http.MethodConnect:
		description = fmt.Sprintf("Get the %s resource connect", tag)
	case http.MethodTrace:
		description = fmt.Sprintf("Get the %s resource trace", tag)
	default:
		description = fmt.Sprintf("Get the %s resource", tag)
	}

	return description
}

// parseRoutePath parses the route path and returns the path and the query parameters.
// Example : /item/:user/:id -> /item/{user}/{id}
func parseRoutePath(route string) (string, []string) {
	// Define a regular expression to match route parameters
	re := regexp.MustCompile(`:([a-zA-Z0-9_]+)`)

	// Find all matches of the parameters in the route
	matches := re.FindAllStringSubmatch(route, -1)

	var params []string

	for _, match := range matches {
		if len(match) > 1 {
			params = append(params, match[1])
		}
	}

	// Replace all instances of :param with {param}
	modifiedRoute := re.ReplaceAllString(route, `{$1}`)

	return modifiedRoute, params
}

// Get status code from the method
func getStatusCode(method string) int {
	switch method {
	case http.MethodPost:
		return http.StatusCreated
	case http.MethodDelete:
		return http.StatusNoContent
	case http.MethodGet, http.MethodPatch, http.MethodPut:
		fallthrough
	default:
		return http.StatusOK
	}
}

// toTitle transform string to title case
func toTitle(s string) string {
	caser := cases.Title(language.Und)

	return caser.String(s)
}
