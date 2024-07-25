package lite

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestDefaultOpenAPIHandler(t *testing.T) {
	app := fiber.New()

	specURL := "http://example.com/swagger.json"
	handler := defaultOpenAPIHandler(specURL)

	app.Get("/*", handler)

	t.Run("returns index.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/html", resp.Header.Get("Content-Type"))
	})

	t.Run("returns root path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/html", resp.Header.Get("Content-Type"))
	})

	t.Run("returns 404 for other paths", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/some-other-path", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
