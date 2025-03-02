package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupOnboardRoutes(app *fiber.App) {
	onboard := app.Group("/onboard", middleware.CheckOnBoardToken())
	onboard.Post("/team/create", handler.CreateTeam)
	onboard.Post("/team/join", handler.JoinTeam)
	onboard.Get("/team/invite", handler.JoinWithInvite)
}