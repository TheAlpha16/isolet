package middleware

import (
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

func IPMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx := c.UserContext()

		ip := c.Get(utils.IPHeaderKey, c.IP())

		// set ip in the context
		userCtx = common.SetFieldInExtraData(userCtx, utils.ContextKeyIP, ip)

		c.SetUserContext(userCtx)
		return c.Next()
	}
}
