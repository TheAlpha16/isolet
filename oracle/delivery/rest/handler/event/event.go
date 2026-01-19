package event

import (
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/response"
	eventDom "github.com/TheAlpha16/isolet/oracle/internal/domain/event"
	"github.com/gofiber/fiber/v2"
)

type EventHandler interface {
	Info(c *fiber.Ctx) error
}

type eventHandler struct {
	eventUc eventDom.Usecase
}

func (h *eventHandler) Info(c *fiber.Ctx) error {
	output, err := h.eventUc.Info(c.UserContext())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("", output))
}

func New(eventUc eventDom.Usecase) EventHandler {
	return &eventHandler{
		eventUc: eventUc,
	}
}
