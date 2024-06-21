package lite

import (
	"github.com/disco07/lite/errors"
	"github.com/disco07/lite/mime"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

type HandlerTestSuite struct {
	suite.Suite
}

func (suite *HandlerTestSuite) SetupTest() {
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestContextWithRequest_ApplicationJSON_Requests() {
	type request struct {
		ID uint64 `lite:"path=id"`
	}

	type response struct {
		ID      uint64 `json:"id"`
		Message string `json:"message"`
	}

	app := New()
	Get(app, "/foo/:id", func(c *ContextWithRequest[request]) (response, error) {
		req, err := c.Requests()
		if err != nil {
			return response{}, err
		}

		return response{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo/123", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

func (suite *HandlerTestSuite) TestContextWithRequest_ApplicationXML_Requests() {
	type request struct {
		ID uint64 `lite:"path=id"`
	}

	type response struct {
		ID      uint64 `xml:"id"`
		Message string `xml:"message"`
	}

	app := New()
	Get(app, "/foo/:id", func(c *ContextWithRequest[request]) (response, error) {
		req, err := c.Requests()
		if err != nil {
			return response{}, err
		}

		c.Type(mime.ApplicationXML)

		return response{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo/123", nil)
	req.Header.Set("Content-Type", "application/xml")
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), "<response><id>123</id><message>Hello World</message></response>", utils.UnsafeString(body))
}

func (suite *HandlerTestSuite) TestContextWithRequest_ApplicationJSON_Requests_Error() {
	type request struct {
		ID uint64 `lite:"path=id"`
	}

	type response struct {
		ID      uint64 `json:"id"`
		Message string `json:"message"`
	}

	app := New()
	Get(app, "/foo/:id", func(c *ContextWithRequest[request]) (response, error) {
		return response{}, assert.AnError
	})

	req := httptest.NewRequest("GET", "/foo/123", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Path_Error() {
	type request struct {
		ID uint64 `lite:"path=id"`
	}

	type response struct {
		ID      uint64 `json:"id"`
		Message string `json:"message"`
	}

	app := New()
	Get(app, "/foo/:id", func(c *ContextWithRequest[request]) (response, error) {
		req, err := c.Requests()
		if err != nil {
			return response{}, err
		}

		if req.ID == 0 {
			return response{}, errors.NewBadRequestError("ID is required")
		}

		return response{}, nil
	})

	req := httptest.NewRequest("GET", "/foo/abc", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 400")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Query() {
	type request struct {
		ID uint64 `lite:"query=id"`
	}

	type response struct {
		ID      uint64 `json:"id"`
		Message string `json:"message"`
	}

	app := New()
	Get(app, "/foo", func(c *ContextWithRequest[request]) (response, error) {
		req, err := c.Requests()
		if err != nil {
			return response{}, err
		}

		return response{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo?id=123", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

func (suite *HandlerTestSuite) TestContextWithRequest_Query_Error() {
	type request struct {
		ID uint64 `lite:"query=id"`
	}

	type response struct {
		ID      uint64 `json:"id"`
		Message string `json:"message"`
	}

	app := New()
	Get(app, "/foo", func(c *ContextWithRequest[request]) (response, error) {
		return response{}, assert.AnError
	})

	req := httptest.NewRequest("GET", "/foo?id=abc", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Header() {
	type request struct {
		ID uint64 `lite:"header=id"`
	}

	type response struct {
		ID      uint64 `json:"id"`
		Message string `json:"message"`
	}

	app := New()
	Get(app, "/foo", func(c *ContextWithRequest[request]) (response, error) {
		req, err := c.Requests()
		if err != nil {
			return response{}, err
		}

		return response{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo", nil)
	req.Header.Set("id", "123")
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

func (suite *HandlerTestSuite) TestContextWithRequest_Header_Error() {
	type request struct {
		ID uint64 `lite:"header=id"`
	}

	type response struct {
		ID      uint64 `json:"id"`
		Message string `json:"message"`
	}

	app := New()
	Get(app, "/foo", func(c *ContextWithRequest[request]) (response, error) {
		return response{}, assert.AnError
	})

	req := httptest.NewRequest("GET", "/foo", nil)
	req.Header.Set("id", "abc")
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type reqBody struct {
	ID float64 `json:"id" xml:"id"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body() {
	type request struct {
		Body reqBody `lite:"req=body"`
	}

	type response struct {
		ID      float64 `json:"id"`
		Message string  `json:"message"`
	}

	app := New()
	Post(app, "/foo", func(c *ContextWithRequest[request]) (response, error) {
		req, err := c.Requests()
		if err != nil {
			return response{}, err
		}

		return response{
			ID:      req.Body.ID,
			Message: "Hello World",
		}, nil
	})

	requestBody := `{"id":123}`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 200")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_Error() {
	type request struct {
		Body reqBody `lite:"req=body"`
	}

	type response struct {
		ID      float64 `json:"id"`
		Message string  `json:"message"`
	}

	app := New()
	Post(app, "/foo", func(c *ContextWithRequest[request]) (response, error) {
		return response{}, assert.AnError
	})

	requestBody := `{"id":"abc"}`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_ApplicationXML() {
	type request struct {
		Body reqBody `lite:"req=body"`
	}

	type response struct {
		ID      float64 `json:"id"`
		Message string  `json:"message"`
	}

	app := New()
	Post(app, "/foo", func(c *ContextWithRequest[request]) (response, error) {
		req, err := c.Requests()
		if err != nil {
			return response{}, err
		}

		return response{
			ID:      req.Body.ID,
			Message: "Hello World",
		}, nil
	})

	requestBody := `<request><id>123</id></request>`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/xml")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 200")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_ApplicationXML_Error() {
	type request struct {
		Body reqBody `lite:"req=body"`
	}

	type response struct {
		ID      float64 `json:"id"`
		Message string  `json:"message"`
	}

	app := New()
	Post(app, "/foo", func(c *ContextWithRequest[request]) (response, error) {
		return response{}, assert.AnError
	})

	requestBody := `<request><id>abc</id></request>`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/xml")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_ApplicationXML_Invalid() {
	type request struct {
		Body reqBody `lite:"req=body"`
	}

	type response struct {
		ID      float64 `json:"id"`
		Message string  `json:"message"`
	}

	app := New()
	Post(app, "/foo", func(c *ContextWithRequest[request]) (response, error) {
		return response{}, assert.AnError
	})

	requestBody := `<request><id>abc</id>`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/xml")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_Put() {
	type request struct {
		Body reqBody `lite:"req=body"`
	}

	type response struct {
		ID      float64 `json:"id"`
		Message string  `json:"message"`
	}

	app := New()
	Put(app, "/foo", func(c *ContextWithRequest[request]) (response, error) {
		req, err := c.Requests()
		if err != nil {
			return response{}, err
		}

		return response{
			ID:      req.Body.ID,
			Message: "Hello World",
		}, nil
	})

	requestBody := `{"id":123}`
	req := httptest.NewRequest("PUT", "/foo", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

func (suite *HandlerTestSuite) TestContextWithRequest_Delete() {
	type request struct {
		ID uint64 `lite:"path=id"`
	}

	app := New()
	Delete(app, "/foo/:id", func(c *ContextWithRequest[request]) (ret struct{}, err error) {
		req, err := c.Requests()
		if err != nil {
			return
		}

		if req.ID == 0 {
			err = errors.NewBadRequestError("ID is required")

			return
		}

		return
	})

	req := httptest.NewRequest("DELETE", "/foo/123", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 204, resp.StatusCode, "Expected status code 204")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Patch() {
	type request struct {
		ID uint64 `lite:"path=id"`
	}

	app := New()
	Patch(app, "/foo/:id", func(c *ContextWithRequest[request]) (ret struct{}, err error) {
		req, err := c.Requests()
		if err != nil {
			return
		}

		if req.ID == 0 {
			err = errors.NewBadRequestError("ID is required")

			return
		}

		return
	})

	req := httptest.NewRequest("PATCH", "/foo/123", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Head() {
	app := New()
	Head(app, "/foo", func(c *ContextNoRequest) (ret struct{}, err error) {

		return
	})

	req := httptest.NewRequest("HEAD", "/foo", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Options() {
	app := New()
	Options(app, "/foo/", func(c *ContextNoRequest) (ret struct{}, err error) {

		return
	})

	req := httptest.NewRequest("OPTIONS", "/foo", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Trace() {
	app := New()
	Trace(app, "/foo/", func(c *ContextNoRequest) (ret struct{}, err error) {

		return
	})

	req := httptest.NewRequest("TRACE", "/foo", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Connect() {
	app := New()
	Connect(app, "/foo/", func(c *ContextNoRequest) (ret struct{}, err error) {

		return
	})

	req := httptest.NewRequest("CONNECT", "/foo", nil)
	resp, err := app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

type fakeContext[r any] struct {
	*ContextNoRequest
}

func (c fakeContext[r]) Requests() (r, error) {
	var ret r
	return ret, nil
}

func (c fakeContext[r]) Status(status int) Context[r] {
	return nil
}

func (c fakeContext[r]) Type(extension string, charset ...string) Context[r] {
	return nil
}

func (suite *HandlerTestSuite) TestCustomContext() {
	assert.Panicsf(suite.T(), func() {
		c := newLiteContext[request, fakeContext[request]](ContextNoRequest{})
		assert.Nil(suite.T(), c)
	}, "unknown type")
}

func (suite *HandlerTestSuite) TestRegisterRoute() {
	app := New()
	type request struct {
		ID uint64 `lite:"path=id"`
	}
	type response struct {
		ID uint64 `json:"id"`
	}

	assert.Panicsf(suite.T(), func() {
		_ = registerRoute[request, response](
			app,
			Route[request, response](Route[response, request]{
				path:        "/foo/:id",
				method:      "GET",
				contentType: "application/json",
				statusCode:  200,
			}),
			nil,
			func(ctx *fiber.Ctx) error {
				return nil
			},
		)
	}, "unknown parameter type")

}
