package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group(AUTH_GROUP, RateLimiter())
	auth.Post(LOGIN, handler.Login)
	auth.Post(REGISTER, middleware.AreRegsOpen(), handler.Register)
	auth.Get(VERIFY, handler.Verify)
	auth.Post(FORGOT_PASSWORD, handler.ForgotPassword)
	auth.Post(RESET_PASSWORD, handler.ResetPassword)
	
	// need to change the group of this route
	auth.Get(METADATA, handler.GetMetadata)
}