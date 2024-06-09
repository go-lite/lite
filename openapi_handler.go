package lite

import "github.com/gofiber/fiber/v2"

func DefaultOpenAPIHandler(specURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(`<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8" />
	<meta name="referrer" content="same-origin" />
	<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	<title>OpenAPI specification</title>
	<script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
	<link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css" />
</head>
<body style="height: 100vh;">
	<elements-api
		apiDescriptionUrl="` + specURL + `"
		layout="responsive"
		router="hash"
		tryItCredentialsPolicy="same-origin"
	/>
</body>
</html>`)
	}
}
