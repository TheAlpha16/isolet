package handler

import (
	"context"

	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	"github.com/gofiber/fiber/v2"
)

type Input interface {
	Validate(ctx context.Context) error
}

func DecodeInput(c *fiber.Ctx, input Input) error {
	if err := c.BodyParser(input); err != nil {
		return errorDom.Raise(c.UserContext(), errorDom.ErrRestInvalidPayload, "", err, nil)
	}

	if err := input.Validate(c.UserContext()); err != nil {
		return err
	}

	return nil
}
