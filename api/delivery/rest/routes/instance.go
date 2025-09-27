package routes

import (
	instanceHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/instance"
	"github.com/TheAlpha16/isolet/api/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/api/infra/jwt"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterInstance(router fiber.Router, instanceHandler instanceHan.InstanceHandler, tokenUc tokenDom.Usecase, jwtSvc jwt.JWT) {
	instanceRouter := router.Group(utils.RouteInstance, middleware.AuthMiddleware(tokenUc, jwtSvc), middleware.RequireTeamMiddleware())
	instanceRouter.Post(utils.RouteInstanceStart, instanceHandler.Start)
}
