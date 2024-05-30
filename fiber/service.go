package fiber

import (
	"context"
	"fmt"
	"github.com/disco07/common-pkg/fiber/service"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"log/slog"
	"strings"
)

func NewHTTPServer(lc fx.Lifecycle, config service.Config, logger *slog.Logger, tracer trace.TracerProvider) fiber.Router {
	svr := fiber.New()
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infof("Starting service 🚀 on port %s\n", config.Port())

			go func() {
				err := svr.Listen(fmt.Sprintf(":%d", config.Port()))
				if err != nil {
					log.Fatal("Error starting server: ", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return svr.Shutdown()
		},
	})

	return router(logger, svr, config, tracer)
}

func router(
	log *slog.Logger,
	f *fiber.App,
	config service.Config,
	tracer trace.TracerProvider,
) *fiber.App {
	f.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOriginsFunc: nil,
		AllowOrigins:     "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowHeaders:     "Accept, Origin, Content-Type, Authorization, X-CSRF-Token",
		ExposeHeaders:    "Link",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	f.Use(otelfiber.Middleware(
		otelfiber.WithPort(config.Port()),
		otelfiber.WithTracerProvider(tracer),
	))

	f.Use(func(c *fiber.Ctx) error {
		log.InfoContext(
			c.Context(),
			"request made",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
		)

		err := c.Next()
		if err != nil {
			log.ErrorContext(
				c.Context(),
				fmt.Sprintf("error while handling request: %v", err),
				slog.Any("error", err),
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
			)
		}

		return err
	})

	return f
}
