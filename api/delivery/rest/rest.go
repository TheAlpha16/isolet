package rest

import (
	"fmt"

	authHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/auth"
	eventHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/event"
	healthHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/health"
	teamHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/team"
	"github.com/TheAlpha16/isolet/api/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/api/delivery/rest/routes"
	"github.com/TheAlpha16/isolet/api/infra"
	"github.com/TheAlpha16/isolet/api/internal/usecase"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	"go.uber.org/zap"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
)

func New(usecases *usecase.Usecases, infra *infra.Infra) *fiber.App {
	config := utils.GetConfig()
	app := fiber.New(
		fiber.Config{
			AppName: config.Name,
		},
	)

	// Setup middlewares
	app.Use(middleware.ContextMiddleware())
	app.Use(otelfiber.Middleware(
		otelfiber.WithNext(func(c *fiber.Ctx) bool {
			return c.Path() == utils.RoutePing
		}),
	))
	app.Use(middleware.SentryMiddleware())
	app.Use(middleware.LoggingMiddleware())
	app.Use(middleware.ErrorMiddleware())

	// Init handlers
	healthHandler := healthHan.New()
	authHandler := authHan.New(usecases.Auth)
	teamHandler := teamHan.New(usecases.Team)
	eventHandler := eventHan.New(usecases.Event)

	// Setup routes
	apiRouter := app.Group(config.Rest.APIVersionPrefix)
	routes.RegisterHealth(apiRouter, healthHandler)
	routes.RegisterAuth(apiRouter, authHandler, usecases.Token, infra.JWT)
	routes.RegisterTeam(apiRouter, teamHandler, usecases.Token, infra.JWT)
	routes.RegisterEvent(apiRouter, eventHandler)
	return app
}

func StartServer(app *fiber.App) {
	config := utils.GetConfig()
	serverPort := fmt.Sprintf(":%d", config.Rest.Port)
	err := app.Listen(serverPort)
	if err != nil {
		logger.GetAppLogger().Fatal("Failed to start REST server", zap.String("port", serverPort), zap.Error(err))
	}
}

func Shutdown(app *fiber.App) {
	logger.GetAppLogger().Info("Shutting down REST server...")
	err := app.Shutdown()
	if err != nil {
		logger.GetAppLogger().Info("Shutting down REST server failed", zap.Error(err))
		return
	}
	logger.GetAppLogger().Info("Shutting down REST server successful")
}
