package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/ping", handler.Ping)
	app.Use(middleware.ResolveIP())

	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Get("/verify", handler.Verify)
	auth.Post("/forgot-password", handler.ForgotPassword)
	auth.Post("/reset-password", handler.ResetPassword)

	onboard := app.Group("/onboard", middleware.CheckOnBoardToken())
	onboard.Post("/team/create", handler.CreateTeam)
	onboard.Post("/team/join", handler.JoinTeam)

	api := app.Group("/api", middleware.CheckTime(), middleware.CheckToken())
	api.Get("/challs", handler.GetChalls)
	api.Post("/submit", handler.SubmitFlag)
	api.Post("/hint/unlock", handler.UnlockHint)

	api.Get("/scoreboard", handler.ShowScoreBoard)
	api.Get("scoreboard/top", handler.GetScoreGraph)
	api.Get("/identify", handler.Identify)
	api.Get("/logout", handler.Logout)

	api.Post("/launch", handler.StartInstance)
	api.Post("/stop", handler.StopInstance)
	api.Get("/status", handler.GetStatus)
	api.Post("/extend", handler.ExtendTime)
}
