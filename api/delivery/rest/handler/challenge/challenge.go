package challenge

import (
	"github.com/TheAlpha16/isolet/api/delivery/rest/response"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	"github.com/gofiber/fiber/v2"
)

type ChallengeHandler interface {
	List(c *fiber.Ctx) error
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

func New(challengeUc challengeDom.Usecase) ChallengeHandler {
	return &challengeHandler{
		challengeUc: challengeUc,
	}
}
