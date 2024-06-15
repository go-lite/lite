package lite

import (
	"github.com/gofiber/fiber/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http/httptest"
	"testing"
)

type CtxTestSuite struct {
	suite.Suite
}

func (suite *CtxTestSuite) SetupTest() {
}

func TestStringTransformTestSuite(t *testing.T) {
	suite.Run(t, new(CtxTestSuite))
}

type request struct {
	ID uint64 `lite:"path=id"`
}

type response struct {
	ID      uint64 `json:"id"`
	Message string `json:"message"`
}

func (suite *CtxTestSuite) TestContext() {
	app := NewApp()
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
