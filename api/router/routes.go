package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/ping", handler.Ping)

	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Get("/verify", handler.Verify)

	api := app.Group("/api", middleware.CheckToken())
	api.Post("/team/create", handler.CreateTeam)
	api.Post("/team/join", handler.JoinTeam)
	api.Get("/challs", handler.GetChalls)
	// api.Post("/launch", handler.StartInstance)
	// api.Post("/stop", handler.StopInstance)
	// api.Post("/submit", handler.SubmitFlag)
	api.Get("/status", handler.GetStatus)
	// api.Get("/scoreboard", handler.ShowScoreBoard)
	// api.Post("/extend", handler.ExtendTime)
}
