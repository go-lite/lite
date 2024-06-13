package lite

import (
	"context"
	"github.com/disco07/lite/errors"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"net/http"
	"regexp"
)

func newLiteContext[Request any, Contexted Context[Request]](ctx ContextNoRequest) Contexted {
	var c Contexted

	switch any(c).(type) {
	case ContextNoRequest:
		return any(ctx).(Contexted)
	case *ContextNoRequest:
		return any(&ctx).(Contexted)
	case *ContextWithRequest[Request]:
		return any(&ContextWithRequest[Request]{
			ContextNoRequest: ctx,
		}).(Contexted)
	default:
		panic("unknown type")
	}
}

func fiberHandler[ResponseBody, Request any, Contexted Context[Request]](
	controller func(c Contexted) (ResponseBody, error),
	path string,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Context().SetContentType("application/json")

		ctx := newLiteContext[Request, Contexted](ContextNoRequest{ctx: c, path: path})

		response, err := controller(ctx)
		if err != nil {
			c.Status(http.StatusInternalServerError)

			return c.JSON(errors.DefaultErrorResponses[http.StatusInternalServerError].SetMessage(err.Error()))
		}

		c.Status(getStatusCode(c.Method()))

		return serializeResponse(c.Context(), &response)
	}
}

func Get[ResponseBody, Request any, Contexted Context[Request]](
	app *App,
	path string,
	controller func(Contexted) (ResponseBody, error),
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

func Post[ResponseBody, Request any, Contexted Context[Request]](
	app *App,
	path string,
	controller func(Contexted) (ResponseBody, error),
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

func Put[ResponseBody, Request any, Contexted Context[Request]](
	app *App,
	path string,
	controller func(Contexted) (ResponseBody, error),
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

func Delete[ResponseBody, Request any, Contexted Context[Request]](
	app *App,
	path string,
	controller func(Contexted) (ResponseBody, error),
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

func Patch[ResponseBody, Request any, Contexted Context[Request]](
	app *App,
	path string,
	controller func(Contexted) (ResponseBody, error),
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

func Head[ResponseBody, Request any, Contexted Context[Request]](
	app *App,
	path string,
	controller func(Contexted) (ResponseBody, error),
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

func Options[ResponseBody, Request any, Contexted Context[Request]](
	app *App,
	path string,
	controller func(Contexted) (ResponseBody, error),
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
	if len(middleware) > 0 {
		app.Add(route.method,
			route.path,
			middleware...,
		)
	}

	app.Add(
		route.method,
		route.path,
		controller,
	)

	operation, err := registerOpenAPIOperation[ResponseBody, Request](
		app,
		route.method,
		route.path,
		route.contentType,
		route.statusCode,
	)
	if err != nil {
		slog.ErrorContext(context.Background(), "failed to register openapi operation", slog.Any("error", err))
		panic(err)
	}

	route.operation = operation

	return route
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
	case http.MethodGet:
		return http.StatusOK
	case http.MethodPost:
		return http.StatusCreated
	case http.MethodPut:
		return http.StatusOK
	case http.MethodDelete:
		return http.StatusNoContent
	case http.MethodPatch:
		return http.StatusOK
	default:
		return http.StatusOK
	}
}
