package middleware

import (
	"github.com/TheAlpha16/isolet/api/infra/jwt"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(tokenUc tokenDom.Usecase, jwtSvc jwt.JWT) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx := c.UserContext()

		token := c.Cookies(utils.AuthTokenCookieName)
		if token == "" {
			return errorDom.Raise(userCtx, errorDom.ErrAuthMissingToken, "", nil, nil)
		}

		claims, err := jwtSvc.Verify(userCtx, token)
		if err != nil {
			return err
		}

		if claims.Purpose != tokenDom.TokenAuth {
			return errorDom.Raise(userCtx, errorDom.ErrTokenExpiredInvalid, "", nil, nil)
		}

		_, err = tokenUc.Fetch(userCtx, &tokenDom.TokenIdentifier{
			ID:       claims.JWTID,
			EntityID: claims.Subject,
			Purpose:  claims.Purpose,
		})
		if err != nil {
			return err
		}

		// set user_id, team_id and role in the context
		userCtx = common.SetFieldInExtraData(userCtx, utils.ContextKeyUserID, *claims.UserID)
		userCtx = common.SetFieldInExtraData(userCtx, utils.ContextKeyRole, *claims.Role)

		if claims.TeamID != nil {
			userCtx = common.SetFieldInExtraData(userCtx, utils.ContextKeyTeamID, *claims.TeamID)
		}

		c.SetUserContext(userCtx)
		return c.Next()
	}
}
