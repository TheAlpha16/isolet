package challenge

import "github.com/gofiber/fiber/v2"

type ChallengeHandler interface {
	List(c *fiber.Ctx) error
}

type challengeHandler struct{}

func (h *challengeHandler) List(c *fiber.Ctx) error {
	return c.SendString("List")
}

func New() ChallengeHandler {
	return &challengeHandler{}
}
