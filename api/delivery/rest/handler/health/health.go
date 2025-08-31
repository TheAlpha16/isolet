package health

import (
	"github.com/TheAlpha16/isolet/api/delivery/rest/response"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler interface {
	CheckHealth(ctx *fiber.Ctx) error
}

type healthHandler struct{}

func (h *healthHandler) CheckHealth(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).
		JSON(response.Success[any]("still alive, atleast for now!", nil))
}

func New() HealthHandler {
	return &healthHandler{}
}
