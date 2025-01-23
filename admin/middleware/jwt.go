package middleware

import (
	"strconv"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/admin/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	jwtware "github.com/gofiber/contrib/jwt"
)

func CheckTime() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime, err := strconv.ParseInt(config.EVENT_START, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "failure",
				"message": "invalid event start time",
			})
		}

		endTime, err := strconv.ParseInt(config.EVENT_END, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "failure",
				"message": "invalid event end time",
			})
		}

		if time.Now().Before(time.Unix(startTime, 0)) {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "failure",
				"message": "event has not started yet",
			})
		}

		if time.Now().After(time.Unix(endTime, 0)) && config.POST_EVENT == "false" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "failure",
				"message": "event has ended",
			})
		}

		return c.Next()
	}
}

func CheckToken() fiber.Handler {

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key:    []byte(config.SESSION_SECRET),
			JWTAlg: jwtware.HS256,
		},
		TokenLookup: "cookie:token",

		SuccessHandler: func(c *fiber.Ctx) error {

			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			userid := int64(claims["userid"].(float64))
			teamid := int64(claims["teamid"].(float64))
			var TeamNameKey models.TeamNameKey

			if !database.UserExists(c, userid) {
				c.ClearCookie("token")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "user does not exist"})
			}

			if teamname, exists := database.TeamExists(c, teamid); exists {
				c.Locals(TeamNameKey, teamname)
				return c.Next()
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "team does not exist"})
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
		TokenLookup: "cookie:token",

		// let this pass if and only if jwt consists of teamid = -1
		SuccessHandler: func(c *fiber.Ctx) error {

			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			userid := int64(claims["userid"].(float64))
			teamid := int64(claims["teamid"].(float64))

			if !database.UserExists(c, userid) {
				c.ClearCookie("token")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "user does not exist"})
			}

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
		"userid":   user.UserID,
		"email":    user.Email,
		"username": user.Username,
		"rank":     user.Rank,
		"teamid":   user.TeamID,
		"exp":      time.Now().Add(time.Hour * time.Duration(config.SESSION_EXP)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.SESSION_SECRET))
	if err != nil {
		return "", err
	}
	return t, nil
}

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
