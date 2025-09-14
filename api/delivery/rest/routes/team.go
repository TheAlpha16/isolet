package routes

import (
	teamHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/team"
	"github.com/TheAlpha16/isolet/api/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/api/infra/jwt"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterTeam(router fiber.Router, teamHandler teamHan.TeamHandler, tokenUc tokenDom.Usecase, jwtSvc jwt.JWT) {
	teamRouter := router.Group(utils.RouteTeam, middleware.AuthMiddleware(tokenUc, jwtSvc))
	teamRouter.Post(utils.RouteTeamCreate, middleware.RequireEmptyTeamMiddleware(), teamHandler.Create)
	teamRouter.Post(utils.RouteTeamJoin, middleware.RequireEmptyTeamMiddleware(), teamHandler.Join)
	teamRouter.Post(utils.RouteTeamInvite, middleware.RequireTeamMiddleware(), middleware.RoleMiddleware([]userDom.Role{userDom.RoleCaptain}), teamHandler.GenerateInvite)
	teamRouter.Get(utils.RouteTeamInvite, middleware.RequireEmptyTeamMiddleware(), teamHandler.AcceptInvite)
}
