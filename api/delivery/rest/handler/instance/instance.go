package instance

import (
	"github.com/TheAlpha16/isolet/api/delivery/rest/handler"
	"github.com/TheAlpha16/isolet/api/delivery/rest/response"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	"github.com/gofiber/fiber/v2"
)

type InstanceHandler interface {
	Start(c *fiber.Ctx) error
}

type instanceHandler struct {
	instanceUc instanceDom.Usecase
}

func (h *instanceHandler) Start(c *fiber.Ctx) error {
	var input instanceDom.StartInput

	if err := handler.DecodeInput(c, &input); err != nil {
		return err
	}

	output, err := h.instanceUc.Start(c.UserContext(), &input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("instance started successfully", output))
}

func New(instanceUc instanceDom.Usecase) InstanceHandler {
	return &instanceHandler{
		instanceUc: instanceUc,
	}
}
