package middleware

import (
	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	"github.com/TheAlpha16/isolet/oracle/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireTeamMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx := c.UserContext()

		if teamID := common.GetFieldFromExtraData[int64](userCtx, utils.ContextKeyTeamID); teamID != 0 {
			return c.Next()
		}

		return errorDom.Raise(userCtx, errorDom.ErrTeamRequired, "", nil, nil)
	}
}

func RequireEmptyTeamMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx := c.UserContext()

		if teamID := common.GetFieldFromExtraData[int64](userCtx, utils.ContextKeyTeamID); teamID == 0 {
			return c.Next()
		}

		return errorDom.Raise(userCtx, errorDom.ErrTeamAlreadyJoined, "", nil, nil)
	}
}
