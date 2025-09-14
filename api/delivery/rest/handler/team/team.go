package team

import "github.com/gofiber/fiber/v2"

type TeamHandler interface {
	Create(c *fiber.Ctx) error
	Join(c *fiber.Ctx) error
}

type teamHandler struct {
}

func (h *teamHandler) Create(c *fiber.Ctx) error {
	return nil
}

func (h *teamHandler) Join(c *fiber.Ctx) error {
	return nil
}

func NewTeamHandler() TeamHandler {
	return &teamHandler{}
}
