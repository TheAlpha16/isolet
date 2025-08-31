package routes

import (
	healthHan "github.com/TheAlpha16/isolet/api/delivery/rest/handler/health"

	"github.com/gofiber/fiber/v2"
)

func RegisterHealth(router fiber.Router, healthHandler healthHan.HealthHandler) {
	router.Get("/ping", healthHandler.CheckHealth)
}
