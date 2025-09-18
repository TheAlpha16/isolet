package score

import (
	"github.com/TheAlpha16/isolet/api/delivery/rest/response"
	scoreDom "github.com/TheAlpha16/isolet/api/internal/domain/score"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

type ScoreHandler interface {
	Scoreboard(c *fiber.Ctx) error
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

func New(scoreUsecase scoreDom.Usecase) ScoreHandler {
	return &scoreHandler{
		scoreUc: scoreUsecase,
	}
}
