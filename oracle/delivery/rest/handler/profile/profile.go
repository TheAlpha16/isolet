package profile

import (
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/response"
	profileDom "github.com/TheAlpha16/isolet/oracle/internal/domain/profile"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler interface {
	Me(c *fiber.Ctx) error
	Team(c *fiber.Ctx) error
}

type profileHandler struct {
	profileUc profileDom.Usecase
}

func (h *profileHandler) Me(c *fiber.Ctx) error {
	output, err := h.profileUc.Me(c.UserContext())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(response.Success("", output))
}

func (h *profileHandler) Team(c *fiber.Ctx) error {
	output, err := h.profileUc.Team(c.UserContext())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(response.Success("", output))
}

func New(profileUc profileDom.Usecase) ProfileHandler {
	return &profileHandler{
		profileUc: profileUc,
	}
}
