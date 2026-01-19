package middleware

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"

	"github.com/gofiber/fiber/v2"
)

// Add extra data and request source to context
func ContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		ctx = common.SetExtraDataInCtx(ctx, make(common.ExtraData))
		ctx = errorDom.SetSrcInCtx(ctx, c.Path()) // Set source to the path
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func ContextDeadlineMiddleware(duration time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), duration)
		defer cancel()
		c.SetUserContext(ctx)

		return c.Next()
	}
}
