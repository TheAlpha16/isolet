package routes

import (
	authHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/auth"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuth(router fiber.Router, authHandler authHan.AuthHandler) {
	authRouter := router.Group("/auth")
	authRouter.Post("/login", authHandler.Login)
	authRouter.Post("/register", authHandler.Register)
}
