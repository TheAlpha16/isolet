package routes

import (
	authHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/auth"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuth(router fiber.Router, authHandler authHan.AuthHandler) {
	authRouter := router.Group(utils.RouteAuth)
	authRouter.Post(utils.RouteAuthLogin, authHandler.Login)
	authRouter.Post(utils.RouteAuthRegister, authHandler.Register)
}
