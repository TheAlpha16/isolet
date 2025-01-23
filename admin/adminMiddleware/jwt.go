package adminmiddleware

import (
	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	jwtware "github.com/gofiber/contrib/jwt"
)

func CheckAdminToken() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.SESSION_SECRET),
			JWTAlg: jwtware.HS256,
		},
		TokenLookup: "cookie:token",

		// admin only
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			userid := int64(claims["userid"].(float64))
			rank := int64(claims["rank"].(float64))
			
			if !database.UserExists(c, userid) {
				c.ClearCookie("token")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "user does not exist"})
			}

			// Check if user is admin (rank == 1)
			if rank == 1 {
				return c.Next()
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "insufficient privileges"})
		},

		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "invalid or expired session token"})
		},
	})
}
