package routes

import (
	challengeHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/challenge"
	"github.com/TheAlpha16/isolet/api/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/api/infra/jwt"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterChallenge(router fiber.Router, challengeHandler challengeHan.ChallengeHandler, tokenUc tokenDom.Usecase, jwtSvc jwt.JWT) {
	challengeRouter := router.Group(utils.RouteChallenge, middleware.AuthMiddleware(tokenUc, jwtSvc), middleware.RequireTeamMiddleware())
	challengeRouter.Get(utils.RouteChallengeList, challengeHandler.List)
	challengeRouter.Post(utils.RouteChallengeSubmitFlag, challengeHandler.SubmitFlag)
	challengeRouter.Post(utils.RouteHintUnlock, challengeHandler.UnlockHint)
}
