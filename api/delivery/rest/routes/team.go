package routes

import (
	teamHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/team"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterTeam(router fiber.Router, teamHandler teamHan.TeamHandler) {
	teamRouter := router.Group(utils.RouteTeam)
	teamRouter.Post(utils.RouteTeamCreate, teamHandler.Create)
	teamRouter.Post(utils.RouteTeamJoin, teamHandler.Join)
}
