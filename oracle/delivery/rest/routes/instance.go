package routes

import (
	instanceHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/instance"
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/middleware"
	"github.com/TheAlpha16/isolet/oracle/infra/jwt"
	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	tokenDom "github.com/TheAlpha16/isolet/oracle/internal/domain/token"
	"github.com/TheAlpha16/isolet/oracle/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterInstance(
	router fiber.Router,
	instanceHandler instanceHan.InstanceHandler,
	tokenUc tokenDom.Usecase,
	jwtSvc jwt.JWT,
	cvUc cvDom.Usecase,
) {
	config := utils.GetConfig()

	instanceRouter := router.Group(
		utils.RouteInstance,
		middleware.AuthMiddleware(tokenUc, jwtSvc),
		middleware.RequireTeamMiddleware(),
		middleware.CheckTimingsMiddleware(cvUc),
	)
	instanceRouter.Get(
		utils.RouteInstanceList,
		instanceHandler.List,
	)
	instanceRouter.Post(
		utils.RouteInstanceStart,
		middleware.ContextDeadlineMiddleware(config.Instances.StartTimeout),
		instanceHandler.Start,
	)
	instanceRouter.Post(
		utils.RouteInstanceStop,
		middleware.ContextDeadlineMiddleware(config.Instances.StopTimeout),
		instanceHandler.Stop,
	)
	instanceRouter.Post(
		utils.RouteInstanceExtend,
		middleware.ContextDeadlineMiddleware(config.Instances.ExtendTimeout),
		instanceHandler.Extend,
	)
}
