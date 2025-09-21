package routes

import (
	authHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/auth"
	"github.com/TheAlpha16/isolet/api/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/api/infra/jwt"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuth(router fiber.Router, authHandler authHan.AuthHandler, tokenUc tokenDom.Usecase, jwtSvc jwt.JWT) {
	authRouter := router.Group(utils.RouteAuth)
	authRouter.Post(utils.RouteAuthLogin, authHandler.Login)
	authRouter.Post(utils.RouteAuthRegister, authHandler.Register)
	authRouter.Get(utils.RouteAuthVerify, authHandler.Verify)
	authRouter.Post(utils.RouteAuthForgotPassword, authHandler.ForgotPassword)
	authRouter.Post(utils.RouteAuthResetPassword, authHandler.ResetPassword)
	authRouter.Get(utils.RouteAuthLogout, middleware.AuthMiddleware(tokenUc, jwtSvc), authHandler.Logout)
}
