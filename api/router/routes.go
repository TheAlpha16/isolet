package router

import (
	"github.com/TitanCrew/isolet/config"
	"github.com/TitanCrew/isolet/handler"
	"github.com/TitanCrew/isolet/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/ping", handler.Ping)

	if !config.DISCORD_FRONTEND {
		auth := app.Group("/auth")
		auth.Post("/login", handler.Login)
		auth.Post("/register", handler.Register)
		auth.Get("/verify", handler.Verify)
		
		api := app.Group("/api", middleware.CheckToken())
		api.Get("/challs", handler.GetChalls)
		api.Post("/launch", handler.StartInstance)
		api.Post("/stop", handler.StopInstance)
		api.Post("/submit", handler.SubmitFlag)
		api.Get("/status", handler.GetStatus)
		api.Get("/scoreboard", handler.ShowScoreBoard)
	} else {
		api := app.Group("/api")
		api.Get("/challs", handler.GetChalls)
		api.Post("/launch", handler.StartInstance)
		api.Post("/stop", handler.StopInstance)
		api.Post("/submit", handler.SubmitFlag)
		api.Get("/status", handler.GetStatus)
		api.Get("/scoreboard", handler.ShowScoreBoard)
	}
}