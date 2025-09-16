package challenge

import (
	"github.com/TheAlpha16/isolet/api/delivery/rest/response"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	"github.com/gofiber/fiber/v2"
)

type ChallengeHandler interface {
	List(c *fiber.Ctx) error
	SubmitFlag(c *fiber.Ctx) error
}

type challengeHandler struct {
	challengeUc challengeDom.Usecase
}

func (h *challengeHandler) List(c *fiber.Ctx) error {
	challenges, err := h.challengeUc.List(c.UserContext())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("", challenges))
}

func (h *challengeHandler) SubmitFlag(c *fiber.Ctx) error {
	input, err := decodeSubmitFlag(c)
	if err != nil {
		return err
	}

	output, err := h.challengeUc.SubmitFlag(c.UserContext(), input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("", output))
}

func New(challengeUc challengeDom.Usecase) ChallengeHandler {
	return &challengeHandler{
		challengeUc: challengeUc,
	}
}
