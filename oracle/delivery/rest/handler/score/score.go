package score

import (
	"github.com/TheAlpha16/isolet/oracle/delivery/rest/response"
	scoreDom "github.com/TheAlpha16/isolet/oracle/internal/domain/score"
	"github.com/TheAlpha16/isolet/oracle/utils"

	"github.com/gofiber/fiber/v2"
)

type ScoreHandler interface {
	Scoreboard(c *fiber.Ctx) error
	ScoreGraph(c *fiber.Ctx) error
}

type scoreHandler struct {
	scoreUc scoreDom.Usecase
}

func (h *scoreHandler) Scoreboard(c *fiber.Ctx) error {
	var input scoreDom.GetScoreboardInput
	input.Page = c.QueryInt(utils.PageQueryKey)
	input.PageSize = c.QueryInt(utils.PageSizeQueryKey)

	if err := input.Validate(c.UserContext()); err != nil {
		return err
	}

	scoreboard, err := h.scoreUc.GetScoreboard(c.UserContext(), &input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("", scoreboard))
}

func (h *scoreHandler) ScoreGraph(c *fiber.Ctx) error {
	scoreGraph, err := h.scoreUc.GetScoreGraph(c.UserContext())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("", scoreGraph))
}

func New(scoreUsecase scoreDom.Usecase) ScoreHandler {
	return &scoreHandler{
		scoreUc: scoreUsecase,
	}
}
