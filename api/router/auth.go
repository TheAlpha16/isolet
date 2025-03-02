package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", middleware.AreRegsOpen(), handler.Register)
	auth.Get("/verify", handler.Verify)
	auth.Post("/forgot-password", handler.ForgotPassword)
	auth.Post("/reset-password", handler.ResetPassword)
	
	// need to change the group of this route
	auth.Get("/metadata", handler.GetMetadata)
}