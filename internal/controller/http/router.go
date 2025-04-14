package router

import (
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	config "github.com/classydevv/fulfillment/configs/providers"
	"github.com/classydevv/fulfillment/internal/controller/http/middleware"
	"github.com/classydevv/fulfillment/internal/controller/http/routes/v1"
	"github.com/classydevv/fulfillment/internal/usecase"
	"github.com/classydevv/fulfillment/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewProviderRouter(app *fiber.App, cfg *config.Config, l logger.Interface, uc usecase.Provider) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("providers")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// K8s probe
	app.Get("/healthz", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(http.StatusOK)
	})

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewProviderRoutes(apiV1Group, uc, l)
	}
}

