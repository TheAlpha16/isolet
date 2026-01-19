package health

import (
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/response"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler interface {
	CheckHealth(c *fiber.Ctx) error
}

type healthHandler struct{}

func (h *healthHandler) CheckHealth(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).
		JSON(response.Success[any]("still alive, atleast for now!", nil))
}

func New() HealthHandler {
	return &healthHandler{}
}
