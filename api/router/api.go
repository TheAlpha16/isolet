package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAPIRoutes(app *fiber.App) {
	api := app.Group(API_GROUP, middleware.CheckToken(), middleware.CheckTime(), RateLimiter())
	api.Get(CHALLS, handler.GetChalls)
	api.Post(SUBMIT, handler.SubmitFlag)
	api.Post(UNLOCK_HINT, handler.UnlockHint)

	scoreboard := api.Group(SCOREBOARD_GROUP)
	scoreboard.Get(SHOW_SCOREBOARD, handler.ShowScoreBoard)
	scoreboard.Get(TOP_SCORES, handler.GetScoreGraph)
	
	api.Get(IDENTIFY, handler.Identify)
	api.Get(LOGOUT, handler.Logout)

	api.Post(LAUNCH, handler.StartInstance)
	api.Post(STOP, handler.StopInstance)
	api.Get(STATUS, handler.GetStatus)
	api.Post(EXTEND, handler.ExtendTime)

	profile := api.Group(PROFILE_GROUP)
	profile.Get(GET_SELF_TEAM, handler.GetSelfTeam)
	profile.Get(GET_INVITE_TOKEN, handler.GetInviteToken)
}