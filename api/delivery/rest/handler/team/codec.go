package team

import (
	"strings"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	"github.com/gofiber/fiber/v2"
)

func decodeCreateRequest(c *fiber.Ctx) (*teamDom.CreateInput, error) {
	var input teamDom.CreateInput
	if err := c.BodyParser(&input); err != nil {
		return nil, errorDom.Raise(c.UserContext(), errorDom.ErrRestInvalidPayload, "", err, nil)
	}

	// normalize
	input.TeamName = strings.TrimSpace(input.TeamName)
	input.Password = strings.TrimSpace(input.Password)

	if err := input.Validate(c.UserContext()); err != nil {
		return nil, err
	}
	return &input, nil
}
