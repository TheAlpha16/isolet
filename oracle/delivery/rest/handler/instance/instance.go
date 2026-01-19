package instance

import (
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/handler"
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/response"
	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	"github.com/gofiber/fiber/v2"
)

type InstanceHandler interface {
	List(c *fiber.Ctx) error
	Start(c *fiber.Ctx) error
	Stop(c *fiber.Ctx) error
	Extend(c *fiber.Ctx) error
}

type instanceHandler struct {
	instanceUc instanceDom.Usecase
}

func (h *instanceHandler) List(c *fiber.Ctx) error {
	instances, err := h.instanceUc.List(c.UserContext())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("", instances))
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

	return c.Status(fiber.StatusCreated).JSON(response.Success("instance started successfully", output))
}

func (h *instanceHandler) Stop(c *fiber.Ctx) error {
	var input instanceDom.StopInput

	if err := handler.DecodeInput(c, &input); err != nil {
		return err
	}

	err := h.instanceUc.Stop(c.UserContext(), &input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success[any]("instance stopped successfully", nil))
}

func (h *instanceHandler) Extend(c *fiber.Ctx) error {
	var input instanceDom.ExtendInput

	if err := handler.DecodeInput(c, &input); err != nil {
		return err
	}

	output, err := h.instanceUc.Extend(c.UserContext(), &input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("instance extended successfully", output))
}

func New(instanceUc instanceDom.Usecase) InstanceHandler {
	return &instanceHandler{
		instanceUc: instanceUc,
	}
}
