package fiber

import (
	"github.com/gofiber/fiber/v2"
	"swagger/codec"
	"swagger/lite"
	"swagger/swagger"
)

type FiberApp struct {
	*fiber.App

	APIDoc *swagger.API
	Prefix string
}

func NewFiberApp(app *fiber.App, apiDoc *swagger.API) FiberApp {
	return FiberApp{
		App:    app,
		APIDoc: apiDoc,
	}
}

func Get[T, B any, E codec.Encoder[T]](app FiberApp, path string, controller func(*lite.ContextWithBody[B]) (T, error), tags []string) fiber.Router {
	return app.Get(path, func(c *fiber.Ctx) error {
		contextWithBody := &lite.ContextWithBody[B]{Ctx: c}
		response, err := controller(contextWithBody)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		var encoder E

		return encoder.Encode(c, response)
	})
}

//func Group(app FiberApp, path string) FiberApp {
//	app.Prefix = path
//
//	return app
//}
//
//func Get[Req, Res any](app FiberApp, path, name string, handler Handler[Req, Res], tags []string) fiber.Router {
//	reqSchema, reqHeaders, reqPathParams, reqQueryParams, reqContentType := generateRequestSchema[Req]()
//	resSchema, resContentType := generateSchema[Res]()
//
//	reqBody := swagger.RequestBody{}
//	if reqContentType != "" {
//		reqBody = swagger.RequestBody{
//			reqContentType: reqSchema,
//		}
//	}
//
//	resBody := swagger.ResponseBody{}
//	if resContentType != "" {
//		resBody = swagger.ResponseBody{
//			200: {resContentType: resSchema.Properties["Value"]},
//		}
//	}
//
//	reqRes := swagger.RequestResponse{
//		Request: swagger.Request{
//			Headers:         reqHeaders,
//			PathParameters:  reqPathParams,
//			QueryParameters: reqQueryParams,
//			Body:            reqBody,
//		},
//		Response: swagger.Response{
//			Body: resBody,
//		},
//	}
//
//	if app.Prefix != "" {
//		path = app.Prefix + path
//	}
//
//	app.APIDoc.RegisterRoute("get", path, name, tags, reqRes)
//	return app.Get(path, customHandler(handler))
//}
//
//func Post[Req, Res any](app FiberApp, path, name string, handler fiber.Handler, tags []string) fiber.Router {
//	reqSchema, reqHeaders, reqPathParams, reqQueryParams, reqContentType := generateRequestSchema[Req]()
//	resSchema, resContentType := generateSchema[Res]()
//
//	reqBody := swagger.RequestBody{}
//	if reqContentType != "" {
//		reqBody = swagger.RequestBody{
//			reqContentType: reqSchema,
//		}
//	}
//
//	resBody := swagger.ResponseBody{}
//	if resContentType != "" {
//		resBody = swagger.ResponseBody{
//			201: {resContentType: resSchema.Properties["Value"]},
//		}
//	}
//
//	reqRes := swagger.RequestResponse{
//		Request: swagger.Request{
//			Headers:         reqHeaders,
//			PathParameters:  reqPathParams,
//			QueryParameters: reqQueryParams,
//			Body:            reqBody,
//		},
//		Response: swagger.Response{
//			Body: resBody,
//		},
//	}
//
//	if app.Prefix != "" {
//		path = app.Prefix + path
//	}
//
//	app.APIDoc.RegisterRoute("post", path, name, tags, reqRes)
//	return app.Post(path, handler)
//}
//
//func Put[Req, Res any](app FiberApp, path, name string, handler fiber.Handler, tags []string) fiber.Router {
//	reqSchema, reqHeaders, reqPathParams, reqQueryParams, reqContentType := generateRequestSchema[Req]()
//	resSchema, resContentType := generateSchema[Res]()
//
//	reqBody := swagger.RequestBody{}
//	if reqContentType != "" {
//		reqBody = swagger.RequestBody{
//			reqContentType: reqSchema,
//		}
//	}
//
//	resBody := swagger.ResponseBody{}
//	if resContentType != "" {
//		resBody = swagger.ResponseBody{
//			200: {resContentType: resSchema.Properties["Value"]},
//		}
//	}
//
//	reqRes := swagger.RequestResponse{
//		Request: swagger.Request{
//			Headers:         reqHeaders,
//			PathParameters:  reqPathParams,
//			QueryParameters: reqQueryParams,
//			Body:            reqBody,
//		},
//		Response: swagger.Response{
//			Body: resBody,
//		},
//	}
//
//	if app.Prefix != "" {
//		path = app.Prefix + path
//	}
//
//	app.APIDoc.RegisterRoute("put", path, name, tags, reqRes)
//	return app.Put(path, handler)
//}
//
//func Patch[Req, Res any](app FiberApp, path, name string, handler fiber.Handler, tags []string) fiber.Router {
//	reqSchema, reqHeaders, reqPathParams, reqQueryParams, reqContentType := generateRequestSchema[Req]()
//	resSchema, resContentType := generateSchema[Res]()
//
//	reqBody := swagger.RequestBody{}
//	if reqContentType != "" {
//		reqBody = swagger.RequestBody{
//			reqContentType: reqSchema,
//		}
//	}
//
//	resBody := swagger.ResponseBody{}
//	if resContentType != "" {
//		resBody = swagger.ResponseBody{
//			200: {resContentType: resSchema.Properties["Value"]},
//		}
//	}
//
//	reqRes := swagger.RequestResponse{
//		Request: swagger.Request{
//			Headers:         reqHeaders,
//			PathParameters:  reqPathParams,
//			QueryParameters: reqQueryParams,
//			Body:            reqBody,
//		},
//		Response: swagger.Response{
//			Body: resBody,
//		},
//	}
//
//	if app.Prefix != "" {
//		path = app.Prefix + path
//	}
//
//	app.APIDoc.RegisterRoute("patch", path, name, tags, reqRes)
//	return app.Patch(path, handler)
//}
//
//func Delete[Req, Res any](app FiberApp, path, name string, handler fiber.Handler, tags []string) fiber.Router {
//	reqSchema, reqHeaders, reqPathParams, reqQueryParams, reqContentType := generateRequestSchema[Req]()
//	resSchema, resContentType := generateSchema[Res]()
//
//	reqBody := swagger.RequestBody{}
//	if reqContentType != "" {
//		reqBody = swagger.RequestBody{
//			reqContentType: reqSchema,
//		}
//	}
//
//	resBody := swagger.ResponseBody{}
//	if resContentType != "" {
//		resBody = swagger.ResponseBody{
//			204: {resContentType: resSchema.Properties["Value"]},
//		}
//	}
//
//	reqRes := swagger.RequestResponse{
//		Request: swagger.Request{
//			Headers:         reqHeaders,
//			PathParameters:  reqPathParams,
//			QueryParameters: reqQueryParams,
//			Body:            reqBody,
//		},
//		Response: swagger.Response{
//			Body: resBody,
//		},
//	}
//
//	if app.Prefix != "" {
//		path = app.Prefix + path
//	}
//
//	app.APIDoc.RegisterRoute("delete", path, name, tags, reqRes)
//	return app.Delete(path, handler)
//}
//
//func Options[Req, Res any](app FiberApp, path, name string, handler fiber.Handler, tags []string) fiber.Router {
//	reqSchema, reqHeaders, reqPathParams, reqQueryParams, reqContentType := generateRequestSchema[Req]()
//	resSchema, resContentType := generateSchema[Res]()
//
//	reqBody := swagger.RequestBody{}
//	if reqContentType != "" {
//		reqBody = swagger.RequestBody{
//			reqContentType: reqSchema,
//		}
//	}
//
//	resBody := swagger.ResponseBody{}
//	if resContentType != "" {
//		resBody = swagger.ResponseBody{
//			200: {resContentType: resSchema.Properties["Value"]},
//		}
//	}
//
//	reqRes := swagger.RequestResponse{
//		Request: swagger.Request{
//			Headers:         reqHeaders,
//			PathParameters:  reqPathParams,
//			QueryParameters: reqQueryParams,
//			Body:            reqBody,
//		},
//		Response: swagger.Response{
//			Body: resBody,
//		},
//	}
//
//	if app.Prefix != "" {
//		path = app.Prefix + path
//	}
//
//	app.APIDoc.RegisterRoute("options", path, name, tags, reqRes)
//	return app.Options(path, handler)
//}
//
//func Head[Req, Res any](app FiberApp, path, name string, handler fiber.Handler, tags []string) fiber.Router {
//	reqSchema, reqHeaders, reqPathParams, reqQueryParams, reqContentType := generateRequestSchema[Req]()
//	resSchema, resContentType := generateSchema[Res]()
//
//	reqBody := swagger.RequestBody{}
//	if reqContentType != "" {
//		reqBody = swagger.RequestBody{
//			reqContentType: reqSchema,
//		}
//	}
//
//	resBody := swagger.ResponseBody{}
//	if resContentType != "" {
//		resBody = swagger.ResponseBody{
//			200: {resContentType: resSchema.Properties["Value"]},
//		}
//	}
//
//	reqRes := swagger.RequestResponse{
//		Request: swagger.Request{
//			Headers:         reqHeaders,
//			PathParameters:  reqPathParams,
//			QueryParameters: reqQueryParams,
//			Body:            reqBody,
//		},
//		Response: swagger.Response{
//			Body: resBody,
//		},
//	}
//
//	if app.Prefix != "" {
//		path = app.Prefix + path
//	}
//
//	app.APIDoc.RegisterRoute("head", path, name, tags, reqRes)
//	return app.Head(path, handler)
//}
//
//func Trace[Req, Res any](app FiberApp, path, name string, handler fiber.Handler, tags []string) fiber.Router {
//	reqSchema, reqHeaders, reqPathParams, reqQueryParams, reqContentType := generateRequestSchema[Req]()
//	resSchema, resContentType := generateSchema[Res]()
//
//	reqBody := swagger.RequestBody{}
//	if reqContentType != "" {
//		reqBody = swagger.RequestBody{
//			reqContentType: reqSchema,
//		}
//	}
//
//	resBody := swagger.ResponseBody{}
//	if resContentType != "" {
//		resBody = swagger.ResponseBody{
//			200: {resContentType: resSchema.Properties["Value"]},
//		}
//	}
//
//	reqRes := swagger.RequestResponse{
//		Request: swagger.Request{
//			Headers:         reqHeaders,
//			PathParameters:  reqPathParams,
//			QueryParameters: reqQueryParams,
//			Body:            reqBody,
//		},
//		Response: swagger.Response{
//			Body: resBody,
//		},
//	}
//
//	if app.Prefix != "" {
//		path = app.Prefix + path
//	}
//
//	app.APIDoc.RegisterRoute("trace", path, name, tags, reqRes)
//	return app.Trace(path, handler)
//}
