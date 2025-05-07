package http

import (
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	config "github.com/classydevv/fulfillment/configs/providers"
	_ "github.com/classydevv/fulfillment/docs" // used by swag tool
	"github.com/classydevv/fulfillment/internal/providers/controller/http/middleware"
	"github.com/classydevv/fulfillment/internal/providers/controller/http/routes/v1"
	"github.com/classydevv/fulfillment/internal/providers/usecase"
	"github.com/classydevv/fulfillment/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Swagger spec:
// @title       Provider API
// @description Service to manager all provider related data: delivery zones and slots, pickup points, tariffs, etc.
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouterProvider(app *fiber.App, cfg *config.Config, l logger.Interface, uc usecase.Provider) {
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
