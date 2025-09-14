package middleware

import (
	"slices"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/gofiber/fiber/v2"
)

// RoleMiddleware checks if the user has any of the required roles.
func RoleMiddleware(roles []userDom.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx := c.UserContext()
		currentRole := common.GetFieldFromExtraData[userDom.Role](userCtx, utils.ContextKeyRole)

		if slices.Contains(roles, currentRole) {
			return c.Next()
		}

		return errorDom.Raise(c.UserContext(), errorDom.ErrUserInsufficentRole, "", nil, nil)
	}
}
