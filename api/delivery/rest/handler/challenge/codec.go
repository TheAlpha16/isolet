package challenge

import (
	"strings"

	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/gofiber/fiber/v2"
)

func decodeSubmitFlag(c *fiber.Ctx) (*challengeDom.SubmitFlagInput, error) {
	var input challengeDom.SubmitFlagInput
	if err := c.BodyParser(&input); err != nil {
		return nil, errorDom.Raise(c.UserContext(), errorDom.ErrRestInvalidPayload, "", err, nil)
	}

	// normalize
	input.Flag = strings.TrimSpace(input.Flag)

	if err := input.Validate(c.UserContext()); err != nil {
		return nil, err
	}

	return &input, nil
}
