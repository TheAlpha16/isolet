package routes

import (
	handler "github.com/TheAlpha16/isolet/admin/handler"
	middleware "github.com/TheAlpha16/isolet/admin/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/ping", handler.Ping)
	app.Use(middleware.ResolveIP())

	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)

	admin := app.Group("/admin")

	// TODO: Authentication removed for testing, fix this
	challenge := admin.Group("/challenges/edit")
	challenge.Post("/data", handler.EditChallengeMetaData)
	challenge.Post("/files", handler.EditChallengeFiles)
	challenge.Post("/requirements", handler.EditChallengeRequirements)
	challenge.Post("/hints", handler.EditChallengeHints)

	config := admin.Group("/config/edit")
	config.Post("/variables", handler.EditConfigValues)
}