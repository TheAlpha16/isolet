package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAPIRoutes(app *fiber.App) {
	api := app.Group(API_GROUP, middleware.CheckToken())
	api.Get(CHALLS, middleware.CheckTime(), RateLimiter(), handler.GetChalls)
	api.Post(SUBMIT, middleware.CheckTime(), RateLimiter(), handler.SubmitFlag)
	api.Post(UNLOCK_HINT, middleware.CheckTime(), RateLimiter(), handler.UnlockHint)

	scoreboard := api.Group(SCOREBOARD_GROUP, RateLimiter())
	scoreboard.Get(SHOW_SCOREBOARD, handler.ShowScoreBoard)
	scoreboard.Get(TOP_SCORES, handler.GetScoreGraph)
	
	api.Get(IDENTIFY, RateLimiter(), handler.Identify)
	api.Get(LOGOUT, RateLimiter(), handler.Logout)

	api.Post(LAUNCH, middleware.CheckTime(), RateLimiter(), handler.StartInstance)
	api.Post(STOP, middleware.CheckTime(), RateLimiter(), handler.StopInstance)
	api.Get(STATUS, middleware.CheckTime(), RateLimiter(), handler.GetStatus)
	api.Post(EXTEND, middleware.CheckTime(), RateLimiter(), handler.ExtendTime)

	profile := api.Group(PROFILE_GROUP, RateLimiter())
	profile.Get(GET_SELF_TEAM, handler.GetSelfTeam)
	profile.Get(GET_INVITE_TOKEN, handler.GetInviteToken)
}