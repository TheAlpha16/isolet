package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAPIRoutes(app *fiber.App) {
	api := app.Group("/api", middleware.CheckTime(), middleware.CheckToken())
	api.Get("/challs", handler.GetChalls)
	api.Post("/submit", handler.SubmitFlag)
	api.Post("/hint/unlock", handler.UnlockHint)

	scoreboard := api.Group("/scoreboard")
	scoreboard.Get("/", handler.ShowScoreBoard)
	scoreboard.Get("/top", handler.GetScoreGraph)
	
	api.Get("/identify", handler.Identify)
	api.Get("/logout", handler.Logout)

	api.Post("/launch", handler.StartInstance)
	api.Post("/stop", handler.StopInstance)
	api.Get("/status", handler.GetStatus)
	api.Post("/extend", handler.ExtendTime)

	profile := api.Group("/profile")
	profile.Get("/team/self", handler.GetSelfTeam)
	profile.Get("/team/invite", handler.GetInviteToken)
}