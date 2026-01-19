package challenge

import (
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/handler"
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/response"
	challengeDom "github.com/TheAlpha16/isolet/oracle/internal/domain/challenge"
	"github.com/gofiber/fiber/v2"
)

type ChallengeHandler interface {
	List(c *fiber.Ctx) error
	SubmitFlag(c *fiber.Ctx) error
	UnlockHint(c *fiber.Ctx) error
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
	var input challengeDom.SubmitFlagInput

	err := handler.DecodeInput(c, &input)
	if err != nil {
		return err
	}

	output, err := h.challengeUc.SubmitFlag(c.UserContext(), &input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("", output))
}

func (h *challengeHandler) UnlockHint(c *fiber.Ctx) error {
	var input challengeDom.UnlockHintInput

	if err := handler.DecodeInput(c, &input); err != nil {
		return err
	}

	hint, err := h.challengeUc.UnlockHint(c.UserContext(), &input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("", hint))
}

func New(challengeUc challengeDom.Usecase) ChallengeHandler {
	return &challengeHandler{
		challengeUc: challengeUc,
	}
}
