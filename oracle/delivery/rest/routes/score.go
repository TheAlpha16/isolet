package routes

import (
	scoreHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/score"
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/oracle/infra/jwt"
	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	tokenDom "github.com/TheAlpha16/isolet/oracle/internal/domain/token"
	"github.com/TheAlpha16/isolet/oracle/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterScore(
	router fiber.Router,
	scoreHandler scoreHan.ScoreHandler,
	tokenUc tokenDom.Usecase,
	jwtSvc jwt.JWT,
	cvUc cvDom.Usecase,
) {
	scoreRouter := router.Group(
		utils.RouteScore,
		middleware.AuthMiddleware(tokenUc, jwtSvc),
		middleware.RequireTeamMiddleware(),
		middleware.CheckTimingsMiddleware(cvUc),
	)
	scoreRouter.Get(utils.RouteScoreboard, scoreHandler.Scoreboard)
	scoreRouter.Get(utils.RouteScoreGraph, scoreHandler.ScoreGraph)
}
