package routes

import (
	scoreHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/score"
	"github.com/TheAlpha16/isolet/api/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/api/infra/jwt"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/gofiber/fiber/v2"
)

func RegisterScore(router fiber.Router, scoreHandler scoreHan.ScoreHandler, tokenUc tokenDom.Usecase, jwtSvc jwt.JWT) {
	scoreRouter := router.Group(
		utils.RouteScore,
		middleware.AuthMiddleware(tokenUc, jwtSvc),
		middleware.RequireTeamMiddleware(),
	)
	scoreRouter.Get(utils.RouteScoreboard, scoreHandler.Scoreboard)
	scoreRouter.Get(utils.RouteScoreGraph, scoreHandler.ScoreGraph)
}
