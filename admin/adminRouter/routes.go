package adminroutes

import (
	adminhandler "github.com/TheAlpha16/isolet/admin/adminHandler"
	adminmiddleware "github.com/TheAlpha16/isolet/admin/adminMiddleware"
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	admin := app.Group("/admin", adminmiddleware.CheckAdminToken())
	admin.Post("/login", handler.Login)
	admin.Post("/edit/challenges", adminhandler.EditChallenges)
	// admin.Post("/edit/teams", adminHandler.EditTeams)
}