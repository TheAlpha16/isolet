package adminroutes

import (
	adminhandler "github.com/TheAlpha16/isolet/admin/adminHandler"
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/ping", handler.Ping)
	app.Use(middleware.ResolveIP())

	admin := app.Group("/admin")
	admin.Post("/login", handler.Login)
	admin.Post("/edit/challenges", adminhandler.EditChallenges)
	// admin.Post("/edit/teams", adminHandler.EditTeams)
}