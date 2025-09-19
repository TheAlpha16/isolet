package routes

import (
	profileHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/profile"
	"github.com/TheAlpha16/isolet/api/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/api/infra/jwt"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterProfile(router fiber.Router, profileHandler profileHan.ProfileHandler, tokenUc tokenDom.Usecase, jwtSvc jwt.JWT) {
	profileRouter := router.Group(utils.RouteProfile)
	profileRouter.Get(utils.RouteProfileMe, middleware.AuthMiddleware(tokenUc, jwtSvc), profileHandler.Me)
	profileRouter.Get(utils.RouteProfileTeam, middleware.AuthMiddleware(tokenUc, jwtSvc), middleware.RequireTeamMiddleware(), profileHandler.Team)
}
