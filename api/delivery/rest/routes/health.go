package routes

import (
	healthHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/health"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterHealth(router fiber.Router, healthHandler healthHan.HealthHandler) {
	router.Get(utils.RoutePing, healthHandler.CheckHealth)
}
