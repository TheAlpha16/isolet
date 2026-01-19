package rest

import (
	"fmt"

	authHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/auth"
	challengeHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/challenge"
	eventHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/event"
	healthHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/health"
	instanceHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/instance"
	profileHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/profile"
	scoreHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/score"
	teamHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/team"
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/routes"
	"github.com/TheAlpha16/isolet/oracle/infra"
	"github.com/TheAlpha16/isolet/oracle/internal/usecase"
	"github.com/TheAlpha16/isolet/oracle/utils"
	"github.com/TheAlpha16/isolet/oracle/utils/logger"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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
	app.Use(middleware.IPMiddleware())
	app.Use(middleware.SentryMiddleware())
	app.Use(middleware.LoggingMiddleware())
	app.Use(middleware.ErrorMiddleware())

	// Init handlers
	healthHandler := healthHan.New()
	authHandler := authHan.New(usecases.Auth)
	teamHandler := teamHan.New(usecases.Team)
	eventHandler := eventHan.New(usecases.Event)
	profileHandler := profileHan.New(usecases.Profile)
	challengeHandler := challengeHan.New(usecases.Challenge)
	scoreHandler := scoreHan.New(usecases.Score)
	instanceHandler := instanceHan.New(usecases.Instance)

	// Setup routes
	apiRouter := app.Group(config.Rest.APIVersionPrefix)
	routes.RegisterHealth(apiRouter, healthHandler)
	routes.RegisterAuth(apiRouter, authHandler, usecases.Token, infra.JWT)
	routes.RegisterTeam(apiRouter, teamHandler, usecases.Token, infra.JWT)
	routes.RegisterEvent(apiRouter, eventHandler)
	routes.RegisterProfile(apiRouter, profileHandler, usecases.Token, infra.JWT)
	routes.RegisterChallenge(apiRouter, challengeHandler, usecases.Token, infra.JWT)
	routes.RegisterScore(apiRouter, scoreHandler, usecases.Token, infra.JWT)
	routes.RegisterInstance(apiRouter, instanceHandler, usecases.Token, infra.JWT)

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
