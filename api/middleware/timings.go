package middleware

import (
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/gofiber/fiber/v2"
)

func CheckTime() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime, err := config.GetInt("EVENT_START")
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "failure",
				"message": "invalid event start time",
			})
		}

		endTime, err := config.GetInt("EVENT_END")
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

		postEvent, err := config.GetBool("POST_EVENT")
		if err != nil {
			postEvent = false
		}

		if time.Now().After(time.Unix(endTime, 0)) && !postEvent {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "failure",
				"message": "event has ended",
			})
		}

		return c.Next()
	}
}
