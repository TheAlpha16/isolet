package middleware

import (
	"github.com/TheAlpha16/isolet/api/config"

	"github.com/gofiber/fiber/v2"
)

func AreRegsOpen() fiber.Handler {
	return func(c *fiber.Ctx) error {
		regsOpen, err := config.GetBool("USER_REGISTRATION")
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "failure",
				"message": "unable to fetch registration status, contact admin",
			})
		}

		if !regsOpen {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "failure",
				"message": "sorry, registrations are paused for the moment",
			})
		}

		return c.Next()
	}
}
