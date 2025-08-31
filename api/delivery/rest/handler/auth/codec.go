package auth

import (
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"

	"github.com/gofiber/fiber/v2"
)

func decodeLoginInput(c *fiber.Ctx) (*authDom.LoginInput, error) {
	var input authDom.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return nil, errorDom.Raise(c.UserContext(), errorDom.ErrInvalidPayload, "", err, nil)
	}
	if err := input.Validate(c.UserContext()); err != nil {
		return nil, err
	}
	return &input, nil
}
