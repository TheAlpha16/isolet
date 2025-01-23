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
	auth.Post("/register", handler.Register)
	auth.Get("/verify", handler.Verify)

	admin := app.Group("/admin", middleware.CheckTime(), middleware.CheckAdminToken())
	admin.Post("/login", handler.Login)
	admin.Post("/edit/challenges", handler.EditChallenges)
	// admin.Post("/edit/teams", adminHandler.EditTeams)
}