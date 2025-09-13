package auth

import (
	"strings"

	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"

	"github.com/gofiber/fiber/v2"
)

func decodeLoginInput(c *fiber.Ctx) (*authDom.LoginInput, error) {
	var input authDom.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return nil, errorDom.Raise(c.UserContext(), errorDom.ErrRestInvalidPayload, "", err, nil)
	}

	// normalize
	input.Identifier = strings.TrimSpace(strings.ToLower(input.Identifier))
	input.Password = strings.TrimSpace(input.Password)

	if err := input.Validate(c.UserContext()); err != nil {
		return nil, err
	}
	return &input, nil
}

func decodeRegisterInput(c *fiber.Ctx) (*authDom.RegisterInput, error) {
	var input authDom.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return nil, errorDom.Raise(c.UserContext(), errorDom.ErrRestInvalidPayload, "", err, nil)
	}

	// normalize
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	if err := input.Validate(c.UserContext()); err != nil {
		return nil, err
	}
	return &input, nil
}

func decodeForgotPasswordInput(c *fiber.Ctx) (*authDom.ForgotPasswordInput, error) {
	var input authDom.ForgotPasswordInput
	if err := c.BodyParser(&input); err != nil {
		return nil, errorDom.Raise(c.UserContext(), errorDom.ErrRestInvalidPayload, "", err, nil)
	}

	// normalize
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	if err := input.Validate(c.UserContext()); err != nil {
		return nil, err
	}
	return &input, nil
}

func decodeResetPasswordInput(c *fiber.Ctx) (*authDom.ResetPasswordInput, error) {
	var input authDom.ResetPasswordInput
	if err := c.BodyParser(&input); err != nil {
		return nil, errorDom.Raise(c.UserContext(), errorDom.ErrRestInvalidPayload, "", err, nil)
	}

	// normalize
	input.Token = strings.TrimSpace(input.Token)
	input.Password = strings.TrimSpace(input.Password)

	if err := input.Validate(c.UserContext()); err != nil {
		return nil, err
	}
	return &input, nil
}
