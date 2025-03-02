package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/ping", handler.Ping)
	app.Use(middleware.ResolveIP())

	SetupAuthRoutes(app)
	SetupOnboardRoutes(app)
	SetupAPIRoutes(app)
}
