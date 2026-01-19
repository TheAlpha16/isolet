package routes

import (
	healthHan "github.com/TheAlpha16/isolet/oracle/delivery/rest/handler/health"
	"github.com/TheAlpha16/isolet/oracle/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterHealth(router fiber.Router, healthHandler healthHan.HealthHandler) {
	router.Get(utils.RoutePing, healthHandler.CheckHealth)
}
