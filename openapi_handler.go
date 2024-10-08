package lite

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// defaultOpenAPIHandler serve Swagger UI with the YAML file spec
func defaultOpenAPIHandler(specURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Params("*") == "index.html" || c.Params("*") == "" {
			c.Type("html")

			return c.SendString(indexHTML(specURL))
		}

		return c.SendStatus(http.StatusNotFound)
	}
}

func indexHTML(specURL string) string {
	return `
   <!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta name="description" content="SwaggerUI" />
  <title>SwaggerUI Go-Lite</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js" crossorigin></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '` + specURL + `',
      dom_id: '#swagger-ui',
    });
  };
</script>
</body>
</html>
    `
}
