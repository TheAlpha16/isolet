package rest

import (
	"fmt"

	"github.com/TheAlpha16/isolet/api/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/api/delivery/rest/routes"
	"github.com/TheAlpha16/isolet/api/internal/usecase"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	"go.uber.org/zap"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
)

func New(
	usecases *usecase.Usecases,
) *fiber.App {
	config := utils.GetConfig()
	app := fiber.New(
		fiber.Config{
			AppName: config.Name,
		},
	)

	app.Use(middleware.ContextMiddleware())
	app.Use(otelfiber.Middleware(
		otelfiber.WithNext(func(ctx *fiber.Ctx) bool {
			return ctx.Path() == "/ping"
		}),
	))
	app.Use(middleware.LoggingMiddleware())
	app.Use(middleware.ErrorMiddleware())

	routes.RegisterHealth(app)

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
