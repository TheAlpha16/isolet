package routes

import (
	challengeHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/challenge"
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/oracle/infra/jwt"
	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	tokenDom "github.com/TheAlpha16/isolet/oracle/internal/domain/token"
	"github.com/TheAlpha16/isolet/oracle/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterChallenge(
	router fiber.Router,
	challengeHandler challengeHan.ChallengeHandler,
	tokenUc tokenDom.Usecase,
	jwtSvc jwt.JWT,
	cvUc cvDom.Usecase,
) {
	challengeRouter := router.Group(
		utils.RouteChallenge,
		middleware.AuthMiddleware(tokenUc, jwtSvc),
		middleware.RequireTeamMiddleware(),
		middleware.CheckTimingsMiddleware(cvUc),
	)
	challengeRouter.Get(utils.RouteChallengeList, challengeHandler.List)
	challengeRouter.Post(utils.RouteChallengeSubmitFlag, challengeHandler.SubmitFlag)
	challengeRouter.Post(utils.RouteHintUnlock, challengeHandler.UnlockHint)
}
