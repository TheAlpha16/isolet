package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupOnboardRoutes(app *fiber.App) {
	onboard := app.Group(ONBOARD_GROUP, middleware.CheckOnBoardToken(), RateLimiter())
	onboard.Post(TEAM_CREATE, handler.CreateTeam)
	onboard.Post(TEAM_JOIN, handler.JoinTeam)
	onboard.Get(JOIN_TEAM_INVITE, handler.JoinWithInvite)
}