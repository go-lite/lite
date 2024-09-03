package lite

import (
	"bytes"
	"fmt"
	"github.com/go-lite/lite/mime"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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

func (suite *HandlerTestSuite) TestUse() {
	app := New(SetValidator(validator.New()))
	Use(app, func(c *fiber.Ctx) error {
		c.Set("test", "test")
		return nil
	})

	req := httptest.NewRequest("GET", "/foo", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	assert.Equal(suite.T(), "test", resp.Header.Get("test"), "Expected test header")
}

func (suite *HandlerTestSuite) TestGroup() {
	app := New(SetValidator(validator.New()))

	newApp := Group(app, "/foo/")

	Get(newApp, "/bar", func(c *ContextNoRequest) (ret struct{}, err error) {
		return
	})

	req := httptest.NewRequest("GET", "/foo/bar", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestGroup2() {
	app := New(SetValidator(validator.New()))

	newApp := Group(app, "/")

	Get(newApp, "/bar", func(c *ContextNoRequest) (ret struct{}, err error) {
		return
	})

	req := httptest.NewRequest("GET", "/bar", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

type requestApplicationJSON struct {
	ID uint64 `lite:"params=id"`
}

type responserequestApplicationJSON struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_ApplicationJSON_Requests() {
	app := New(SetValidator(validator.New()))

	Get(app, "/foo/:id", func(c *ContextWithRequest[requestApplicationJSON]) (responserequestApplicationJSON, error) {
		req, err := c.Requests()
		if err != nil {
			return responserequestApplicationJSON{}, err
		}

		return responserequestApplicationJSON{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo/123", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

type requestApplicationXML struct {
	ID uint64 `lite:"params=id"`
}

type responseApplicationXML struct {
	ID      uint64 `xml:"id"`
	Message string `xml:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_ApplicationXML_Requests() {
	app := New(SetValidator(validator.New()))
	Get(app, "/foo/:id", func(c *ContextWithRequest[requestApplicationXML]) (responseApplicationXML, error) {
		req, err := c.Requests()
		if err != nil {
			return responseApplicationXML{}, err
		}

		c.SetContentType(mime.ApplicationXML)

		return responseApplicationXML{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo/123", nil)
	req.Header.Set("Content-Type", "application/xml")
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(suite.T(), "<responseApplicationXML><id>123</id><message>Hello World</message></responseApplicationXML>", utils.UnsafeString(body))
}

type requestApplicationJSONError struct {
	ID uint64 `lite:"params=id"`
}

type responseApplicationJSONError struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_ApplicationJSON_Requests_Error() {
	app := New(SetValidator(validator.New()))
	Get(app, "/foo/:id", func(c *ContextWithRequest[requestApplicationJSONError]) (responseApplicationJSONError, error) {
		return responseApplicationJSONError{}, assert.AnError
	})

	req := httptest.NewRequest("GET", "/foo/123", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type requestPath struct {
	ID uint64 `lite:"params=id"`
}

type responsePath struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Path_Error() {
	app := New(SetValidator(validator.New()))
	Get(app, "/foo/:id", func(c *ContextWithRequest[requestPath]) (responsePath, error) {
		req, err := c.Requests()
		if err != nil {
			return responsePath{}, err
		}

		if req.ID == 0 {
			return responsePath{}, NewBadRequestError("ID is required")
		}

		return responsePath{}, nil
	})

	req := httptest.NewRequest("GET", "/foo/abc", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 400")
}

type requestQuery struct {
	ID uint64 `lite:"query=id"`
}

type responseQuery struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Query() {
	app := New(SetValidator(validator.New()))
	Get(app, "/foo", func(c *ContextWithRequest[requestQuery]) (responseQuery, error) {
		req, err := c.Requests()
		if err != nil {
			return responseQuery{}, err
		}

		return responseQuery{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo?id=123", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

func (suite *HandlerTestSuite) TestContextWithRequest_Query_Error() {
	app := New(SetValidator(validator.New()))
	Get(app, "/foo", func(c *ContextWithRequest[requestPath]) (responsePath, error) {
		return responsePath{}, assert.AnError
	})

	req := httptest.NewRequest("GET", "/foo?id=abc", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type requestHeader struct {
	ID uint64 `lite:"header=id,isauth,type=apiKey,name=id"`
}

type responseHeader struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Header() {
	app := New(SetValidator(validator.New()))
	Get(app, "/foo", func(c *ContextWithRequest[requestHeader]) (responseHeader, error) {
		req, err := c.Requests()
		if err != nil {
			return responseHeader{}, err
		}

		return responseHeader{
			ID:      req.ID,
			Message: "Hello World",
		}, nil
	})

	req := httptest.NewRequest("GET", "/foo", nil)
	req.Header.Set("id", "123")
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

type requestHeaderErr struct {
	ID uint64 `lite:"header=id"`
}

type responseHeaderErr struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Header_Error() {
	app := New(SetValidator(validator.New()))
	Get(app, "/foo", func(c *ContextWithRequest[requestHeaderErr]) (responseHeaderErr, error) {
		return responseHeaderErr{}, assert.AnError
	})

	req := httptest.NewRequest("GET", "/foo", nil)
	req.Header.Set("id", "abc")
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type reqBody struct {
	ID float64 `json:"id" xml:"id"`
}

type requestBody struct {
	Body reqBody `lite:"req=body"`
}

type responseBody struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body() {
	app := New(SetValidator(validator.New()))
	Post(app, "/foo", func(c *ContextWithRequest[requestBody]) (responseBody, error) {
		req, err := c.Requests()
		if err != nil {
			return responseBody{}, err
		}

		return responseBody{
			ID:      req.Body.ID,
			Message: "Hello World",
		}, nil
	})

	bodyJSON := `{"id":123}`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 200")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

type requestBodyError struct {
	Body reqBody `lite:"req=body"`
}

type responseBodyError struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_Error() {
	app := New(SetValidator(validator.New()))
	Post(app, "/foo", func(c *ContextWithRequest[requestBodyError]) (responseBodyError, error) {
		return responseBodyError{}, assert.AnError
	})

	bodyJSON := `{"id":"abc"}`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type requestBodyXML struct {
	Body reqBody `lite:"req=body"`
}

type responseBodyXML struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_ApplicationXML() {
	app := New(SetValidator(validator.New()))
	Post(app, "/foo", func(c *ContextWithRequest[requestBodyXML]) (responseBodyXML, error) {
		req, err := c.Requests()
		if err != nil {
			return responseBodyXML{}, err
		}

		return responseBodyXML{
			ID:      req.Body.ID,
			Message: "Hello World",
		}, nil
	})

	bdyXML := `<request><id>123</id></request>`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(bdyXML))
	req.Header.Set("Content-Type", "application/xml")

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 200")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

type requestBodyXMLError struct {
	Body reqBody `lite:"req=body"`
}

type responseBodyXMLError struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_ApplicationXML_Error() {
	app := New(SetValidator(validator.New()))
	Post(app, "/foo", func(c *ContextWithRequest[requestBodyXMLError]) (responseBodyXMLError, error) {
		return responseBodyXMLError{}, assert.AnError
	})

	bodyXML := `<request><id>abc</id></request>`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(bodyXML))
	req.Header.Set("Content-Type", "application/xml")

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type requestBodyXMLInvalid struct {
	Body reqBody `lite:"req=body"`
}

type responseBodyXMLInvalid struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_ApplicationXML_Invalid() {
	app := New(SetValidator(validator.New()))
	Post(app, "/foo", func(c *ContextWithRequest[requestBodyXMLInvalid]) (responseBodyXMLInvalid, error) {
		return responseBodyXMLInvalid{}, assert.AnError
	})

	requestBody := `<request><id>abc</id>`
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/xml")

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 500, resp.StatusCode, "Expected status code 500")
}

type requestBodyPut struct {
	Body reqBody `lite:"req=body"`
}

type responseBodyPut struct {
	ID      float64 `json:"id"`
	Message string  `json:"message"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_Put() {
	app := New(SetValidator(validator.New()))
	Put(app, "/foo", func(c *ContextWithRequest[requestBodyPut]) (responseBodyPut, error) {
		req, err := c.Requests()
		if err != nil {
			return responseBodyPut{}, err
		}

		return responseBodyPut{
			ID:      req.Body.ID,
			Message: "Hello World",
		}, nil
	})

	bodyJSON := `{"id":123}`
	req := httptest.NewRequest("PUT", "/foo", strings.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"message":"Hello World"}`, utils.UnsafeString(body))
}

type requestDelete struct {
	ID uint64 `lite:"params=id"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Delete() {
	app := New(SetValidator(validator.New()))
	Delete(app, "/foo/:id", func(c *ContextWithRequest[requestDelete]) (ret struct{}, err error) {
		req, err := c.Requests()
		if err != nil {
			return
		}

		if req.ID == 0 {
			err = NewBadRequestError("ID is required")

			return
		}

		return
	})

	req := httptest.NewRequest("DELETE", "/foo/123", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 204, resp.StatusCode, "Expected status code 204")
}

type requestPatch struct {
	ID uint64 `lite:"params=id"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_Patch() {
	app := New(SetValidator(validator.New()))
	Patch(app, "/foo/:id", func(c *ContextWithRequest[requestPatch]) (ret struct{}, err error) {
		req, err := c.Requests()
		if err != nil {
			return
		}

		if req.ID == 0 {
			err = NewBadRequestError("ID is required")

			return
		}

		return
	})

	req := httptest.NewRequest("PATCH", "/foo/123", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

type requestPatchError struct {
	ID uint64 `lite:"params=id"`
}

func (suite *HandlerTestSuite) TestContextWithRequest_PatchError() {
	app := New(SetValidator(validator.New()))
	Patch(app, "/foo/:id", func(c *ContextWithRequest[requestPatchError]) (ret struct{}, err error) {
		req, err := c.Requests()
		if err != nil {
			return
		}

		if req.ID == 0 {
			err = NewBadRequestError("ID is required")

			return
		}

		return
	})

	req := httptest.NewRequest("PATCH", "/foo/0", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 400, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Head() {
	app := New(SetValidator(validator.New()))
	Head(app, "/foo", func(c *ContextNoRequest) (ret struct{}, err error) {
		c.SetContentType(mime.ApplicationJSON)
		c.Status(200)

		return
	})

	req := httptest.NewRequest("HEAD", "/foo", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Options() {
	app := New(SetValidator(validator.New()))
	Options(app, "/foo/", func(c *ContextNoRequest) (ret struct{}, err error) {

		return
	})

	req := httptest.NewRequest("OPTIONS", "/foo", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Trace() {
	app := New(SetValidator(validator.New()))
	Trace(app, "/foo/", func(c *ContextNoRequest) (ret struct{}, err error) {

		return
	})

	req := httptest.NewRequest("TRACE", "/foo", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Connect() {
	app := New(SetValidator(validator.New()))
	Connect(app, "/foo/", func(c *ContextNoRequest) (ret struct{}, err error) {

		return
	})

	req := httptest.NewRequest("CONNECT", "/foo", nil)
	resp, err := app.app.Test(req)
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

func (c fakeContext[r]) SetContentType(extension string, charset ...string) Context[r] {
	return nil
}

func (suite *HandlerTestSuite) TestCustomContext() {
	assert.Panicsf(suite.T(), func() {
		c := newLiteContext[request, fakeContext[request]](ContextNoRequest{})
		assert.Nil(suite.T(), c)
	}, "unknown type")
}

type requestRoute struct {
	ID uint64 `lite:"params=id"`
}

type responseRoute struct {
	ID uint64 `json:"id"`
}

func (suite *HandlerTestSuite) TestRegisterRoute() {
	app := New(SetValidator(validator.New()))

	assert.Panicsf(suite.T(), func() {
		_ = registerRoute[requestRoute, responseRoute](
			app,
			Route[requestRoute, responseRoute]{
				path:        "/foo/:id",
				method:      "GET",
				contentType: "application/json",
				statusCode:  200,
			},
			nil,
			func(ctx *fiber.Ctx) error {
				return nil
			},
		)
	}, "unknown parameter type")

}

func (suite *HandlerTestSuite) TestContextWithRequest_FullBody() {
	app := New(SetValidator(validator.New()))

	Post(app, "/test/:id/:is_admin", func(c *ContextWithRequest[testRequest]) (testResponse, error) {
		req, err := c.Requests()
		if err != nil {
			return testResponse{}, err
		}

		if req.Filter == nil {
			c.Set("filter", "test")
		}

		assert.Equal(suite.T(), "", c.Get("User-Agent"))

		err = c.SaveFile(req.Body.File, "./logo/lite.png")
		if err != nil {
			return testResponse{}, err
		}

		method := c.Method()
		if method != http.MethodPost {
			return testResponse{}, NewBadRequestError("Method is not POST")
		}

		return testResponse{
			ID:        req.Params.ID,
			FirstName: req.Body.Metadata.FirstName,
			LastName:  req.Body.Metadata.LastName,
			Gender:    Male,
			GenderSlice: []Gender{
				Male,
				Female,
			},
		}, nil
	})

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Ajouter le fichier
	fileWriter, err := w.CreateFormFile("file", "lite.png")
	if err != nil {
		suite.T().Fatalf("Failed to create form file: %s", err)
	}

	file, err := os.Open("./logo/lite.png")
	if err != nil {
		suite.T().Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		suite.T().Fatalf("Failed to copy file: %s", err)
	}

	// Ajouter le JSON
	metadataWriter, err := w.CreateFormField("metadata")
	if err != nil {
		suite.T().Fatalf("Failed to create form field: %s", err)
	}
	data := `{"first_name":"John","last_name":"Doe", "birthday": "2000-01-01T00:00:00Z"}`
	_, err = metadataWriter.Write([]byte(data))
	if err != nil {
		suite.T().Fatalf("Failed to write metadata: %s", err)
	}

	nameWriter, err := w.CreateFormField("name")
	if err != nil {
		suite.T().Fatalf("Failed to create form field: %s", err)
	}
	_, err = nameWriter.Write([]byte("test"))
	if err != nil {
		suite.T().Fatalf("Failed to write name: %s", err)
	}

	w.Close()

	req := httptest.NewRequest("POST", "/test/123/true", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 201")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"id":123,"first_name":"John","last_name":"Doe", "gender":"male", "gender_slice":["male","female"], "name":""}`, utils.UnsafeString(body))
}

func (suite *HandlerTestSuite) TestContextWithRequest_MultiFile() {
	app := New(SetValidator(validator.New()))
	Post(app, "/foo", func(c *ContextWithRequest[testRequestMultiFileRequest]) (testResponse, error) {
		req, err := c.Requests()
		if err != nil {
			return testResponse{}, err
		}

		if req.Body.Files == nil {
			return testResponse{}, NewBadRequestError("Files are required")
		}

		return testResponse{}, nil
	})

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Ajouter le fichier
	fileWriter, err := w.CreateFormFile("files", "lite.png")
	if err != nil {
		suite.T().Fatalf("Failed to create form file: %s", err)
	}

	file, err := os.Open("./logo/lite.png")
	if err != nil {
		suite.T().Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		suite.T().Fatalf("Failed to copy file: %s", err)
	}

	w.Close()

	req := httptest.NewRequest("POST", "/foo", strings.NewReader(b.String()))
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 201")

	spec, err := app.saveOpenAPISpec()
	assert.NoError(suite.T(), err)
	fmt.Println(string(spec))
}

type requestBodyApplicationPDF struct {
	Body []byte `lite:"req=body,application/pdf"`
}

type responseBodyApplicationPDF = []byte

func (suite *HandlerTestSuite) TestContextWithRequest_Body_ApplicationPDF() {
	app := New(SetValidator(validator.New()))
	Post(app, "/foo", func(c *ContextWithRequest[requestBodyApplicationPDF]) (responseBodyApplicationPDF, error) {
		req, err := c.Requests()
		if err != nil {
			return responseBodyApplicationPDF{}, err
		}

		c.SetContentType(mime.ApplicationPdf)

		return req.Body, nil
	})

	bodyPDF := `%PDF-1.4
%âãÏÓ
1 0 obj
  << /Type /Catalog
     /Pages 2 0 R
  >>
endobj

2 0 obj
  << /Type /Pages
     /Kids [3 0 R]
     /Count 1
     /MediaBox [0 0 300 144]
  >>
endobj

3 0 obj
  <<  /Type /Page
      /Parent 2 0 R
      /Resources
       << /Font
           << /F1
               << /Type /Font
                  /Subtype /Type1
                  /BaseFont /Times-Roman
               >>
           >>
       >>
      /Contents 4 0 R
   >>
endobj

4 0 obj
  << /Length 55 >>
stream
  BT
    /F1 18 Tf
    0 0 Td
    (Hello World) Tj
  ET
endstream
endobj

xref
0 5
0000000000 65535 f 
0000000018 00000 n 
0000000077 00000 n 
0000000178 00000 n 
0000000457 00000 n 
trailer
  <<  /Root 1 0 R
      /Size 5
  >>
startxref
565
%%EOF`

	req := httptest.NewRequest("POST", "/foo", strings.NewReader(bodyPDF))
	req.Header.Set("Content-Type", "application/pdf")

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 201")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), `%PDF-1.4
%âãÏÓ
1 0 obj
  << /Type /Catalog
     /Pages 2 0 R
  >>
endobj

2 0 obj
  << /Type /Pages
     /Kids [3 0 R]
     /Count 1
     /MediaBox [0 0 300 144]
  >>
endobj

3 0 obj
  <<  /Type /Page
      /Parent 2 0 R
      /Resources
       << /Font
           << /F1
               << /Type /Font
                  /Subtype /Type1
                  /BaseFont /Times-Roman
               >>
           >>
       >>
      /Contents 4 0 R
   >>
endobj

4 0 obj
  << /Length 55 >>
stream
  BT
    /F1 18 Tf
    0 0 Td
    (Hello World) Tj
  ET
endstream
endobj

xref
0 5
0000000000 65535 f 
0000000018 00000 n 
0000000077 00000 n 
0000000178 00000 n 
0000000457 00000 n 
trailer
  <<  /Root 1 0 R
      /Size 5
  >>
startxref
565
%%EOF`, utils.UnsafeString(body))

	spec, err := app.saveOpenAPISpec()
	assert.NoError(suite.T(), err)

	expected := `components:
    schemas:
        Body:
            format: binary
            type: string
        httpGenericError:
            properties:
                '@context':
                    type: string
                '@type':
                    type: string
                description:
                    type: string
                status:
                    type: integer
                title:
                    type: string
                violations:
                    items:
                        properties:
                            code:
                                type: string
                            message:
                                type: string
                            more:
                                additionalProperties: {}
                                type: object
                            propertyPath:
                                type: string
                        type: object
                    type: array
            type: object
        uint8:
            format: binary
            type: string
info:
    description: OpenAPI
    title: OpenAPI
    version: 0.0.1
openapi: 3.0.3
paths:
    /foo:
        post:
            requestBody:
                content:
                    application/pdf:
                        schema:
                            $ref: '#/components/schemas/Body'
            responses:
                "201":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/uint8'
                    description: OK
                "400":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Bad Request
                "500":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Internal Server Error
`

	assert.YAMLEqf(suite.T(), expected, string(spec), "openapi generated spec")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_Slice_Uint8() {
	app := New(SetValidator(validator.New()))
	Post(app, "/foo", func(c *ContextWithRequest[[]byte]) (string, error) {
		req, err := c.Requests()
		if err != nil {
			return "", err
		}

		c.SetContentType(mime.TextPlain)

		return string(req), nil
	}).SetResponseContentType("text/plain")

	req := httptest.NewRequest("POST", "/foo", strings.NewReader("Hello World"))
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 201")
	b, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), `Hello World`, utils.UnsafeString(b))

	spec, err := app.saveOpenAPISpec()
	assert.NoError(suite.T(), err)

	expected := `components:
    schemas:
        httpGenericError:
            properties:
                '@context':
                    type: string
                '@type':
                    type: string
                description:
                    type: string
                status:
                    type: integer
                title:
                    type: string
                violations:
                    items:
                        properties:
                            code:
                                type: string
                            message:
                                type: string
                            more:
                                additionalProperties: {}
                                type: object
                            propertyPath:
                                type: string
                        type: object
                    type: array
            type: object
        string:
            type: string
        uint8:
            format: binary
            type: string
info:
    description: OpenAPI
    title: OpenAPI
    version: 0.0.1
openapi: 3.0.3
paths:
    /foo:
        post:
            requestBody:
                content:
                    application/octet-stream:
                        schema:
                            $ref: '#/components/schemas/uint8'
            responses:
                "201":
                    content:
                        text/plain:
                            schema:
                                $ref: '#/components/schemas/string'
                    description: OK
                "400":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Bad Request
                "500":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Internal Server Error`
	assert.YAMLEqf(suite.T(), expected, string(spec), "openapi generated spec")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_String() {
	app := New(SetValidator(validator.New()))
	Post(app, "/foo", func(c *ContextWithRequest[string]) (string, error) {
		req, err := c.Requests()
		if err != nil {
			return "", err
		}

		c.SetContentType(mime.TextPlain)

		return req, nil
	}).SetResponseContentType("text/plain")

	req := httptest.NewRequest("POST", "/foo", strings.NewReader("Hello World"))
	req.Header.Set("Content-Type", "text/plain")

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 201")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), `Hello World`, utils.UnsafeString(body))

	spec, err := app.saveOpenAPISpec()
	assert.NoError(suite.T(), err)

	expected := `components:
    schemas:
        httpGenericError:
            properties:
                '@context':
                    type: string
                '@type':
                    type: string
                description:
                    type: string
                status:
                    type: integer
                title:
                    type: string
                violations:
                    items:
                        properties:
                            code:
                                type: string
                            message:
                                type: string
                            more:
                                additionalProperties: {}
                                type: object
                            propertyPath:
                                type: string
                        type: object
                    type: array
            type: object
        string:
            type: string
info:
    description: OpenAPI
    title: OpenAPI
    version: 0.0.1
openapi: 3.0.3
paths:
    /foo:
        post:
            requestBody:
                content:
                    text/plain:
                        schema:
                            $ref: '#/components/schemas/string'
            responses:
                "201":
                    content:
                        text/plain:
                            schema:
                                $ref: '#/components/schemas/string'
                    description: OK
                "400":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Bad Request
                "500":
                    content:
                        application/json:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        application/xml:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                        multipart/form-data:
                            schema:
                                $ref: '#/components/schemas/httpGenericError'
                    description: Internal Server Error`
	assert.YAMLEqf(suite.T(), expected, string(spec), "openapi generated spec")
}

func (suite *HandlerTestSuite) TestContextWithRequest_Body_StringJSONSwagger() {
	app := New(
		SetTypeOfExtension(JSONExtension),
	)
	Post(app, "/foo", func(c *ContextWithRequest[string]) (string, error) {
		req, err := c.Requests()
		if err != nil {
			return "", err
		}

		c.SetContentType(mime.TextPlain)

		return req, nil
	}).SetResponseContentType("text/plain")

	req := httptest.NewRequest("POST", "/foo", strings.NewReader("Hello World"))
	req.Header.Set("Content-Type", "text/plain")

	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), 201, resp.StatusCode, "Expected status code 201")
	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err, "Expected no error")
	assert.Equal(suite.T(), `Hello World`, utils.UnsafeString(body))

	spec, err := app.saveOpenAPISpec()
	assert.NoError(suite.T(), err)

	expected := `{"components":{"schemas":{"httpGenericError":{"properties":{"@context":{"type":"string"},"@type":{"type":"string"},"description":{"type":"string"},"status":{"type":"integer"},"title":{"type":"string"},"violations":{"items":{"properties":{"code":{"type":"string"},"message":{"type":"string"},"more":{"additionalProperties":{},"type":"object"},"propertyPath":{"type":"string"}},"type":"object"},"type":"array"}},"type":"object"},"string":{"type":"string"}}},"info":{"description":"OpenAPI","title":"OpenAPI","version":"0.0.1"},"openapi":"3.0.3","paths":{"/foo":{"post":{"requestBody":{"content":{"text/plain":{"schema":{"$ref":"#/components/schemas/string"}}}},"responses":{"201":{"content":{"text/plain":{"schema":{"$ref":"#/components/schemas/string"}}},"description":"OK"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/httpGenericError"}},"application/xml":{"schema":{"$ref":"#/components/schemas/httpGenericError"}},"multipart/form-data":{"schema":{"$ref":"#/components/schemas/httpGenericError"}}},"description":"Bad Request"},"500":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/httpGenericError"}},"application/xml":{"schema":{"$ref":"#/components/schemas/httpGenericError"}},"multipart/form-data":{"schema":{"$ref":"#/components/schemas/httpGenericError"}}},"description":"Internal Server Error"}}}}}}`
	assert.JSONEq(suite.T(), expected, string(spec), "openapi generated spec")
}

func (suite *HandlerTestSuite) TestContextWithRequest_SimpleReturnList() {
	app := New(SetValidator(validator.New()))
	Get(app, "/foo", func(c *ContextNoRequest) (ret List[string], err error) {
		return
	})

	req := httptest.NewRequest("GET", "/foo", nil)
	resp, err := app.app.Test(req)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 200, resp.StatusCode, "Expected status code 200")
}

func (suite *HandlerTestSuite) TestSetDescription() {
	assert.Equal(suite.T(), "Get the Test resource", setDescription(http.MethodGet, toTitle("test")))
	assert.Equal(suite.T(), "Create a new test resource", setDescription(http.MethodPost, "test"))
	assert.Equal(suite.T(), "Replace the test resource", setDescription(http.MethodPut, "test"))
	assert.Equal(suite.T(), "Update the test resource", setDescription(http.MethodPatch, "test"))
	assert.Equal(suite.T(), "Delete the test resource", setDescription(http.MethodDelete, "test"))
	assert.Equal(suite.T(), "Get the test resource header", setDescription(http.MethodHead, "test"))
	assert.Equal(suite.T(), "Get the test resource options", setDescription(http.MethodOptions, "test"))
	assert.Equal(suite.T(), "Get the test resource connect", setDescription(http.MethodConnect, "test"))
	assert.Equal(suite.T(), "Get the test resource trace", setDescription(http.MethodTrace, "test"))
	assert.Equal(suite.T(), "Get the test resource", setDescription("", "test"))
}
