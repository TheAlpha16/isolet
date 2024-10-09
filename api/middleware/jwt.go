package middleware

import (
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CheckToken() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key:    []byte(config.SESSION_SECRET),
			JWTAlg: jwtware.HS256,
		},

		SuccessHandler: func(c *fiber.Ctx) error {

			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			teamid := int(claims["teamid"].(float64))

			if teamid != -1 {
				return c.Next()
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "please create or join a team"})
		},

		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "invalid or expired session token"})
		},
	})
}

func CheckOnBoardToken() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key:    []byte(config.SESSION_SECRET),
			JWTAlg: jwtware.HS256,
		},

		// let this pass if and only if jwt consists of teamid = -1
		SuccessHandler: func(c *fiber.Ctx) error {

			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			teamid := int(claims["teamid"].(float64))

			if teamid == -1 {
				return c.Next()
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "already on a team"})
		},

		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "invalid or expired session token"})
		},
	})
}

func GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"userid": user.UserID,
		"email":  user.Email,
		"rank":   user.Rank,
		"teamid": user.TeamID,
		"exp":    time.Now().Add(time.Hour * time.Duration(config.SESSION_EXP)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.SESSION_SECRET))
	if err != nil {
		return "", err
	}
	return t, nil
}
