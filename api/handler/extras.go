package handler

import (
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/gofiber/fiber/v2"
)

func Ping(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Server is up!"})
}

func GetMetadata(c *fiber.Ctx) error {
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

	postEvent, err := config.GetBool("POST_EVENT")
	if err != nil {
		postEvent = false
	}

	teamLen, err := config.GetInt("TEAM_LEN")
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "failure",
			"message": "invalid team length",
		})
	}

	ctfName, err := config.Get("CTF_NAME")
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "failure",
			"message": "invalid CTF name",
		})
	}

	return c.Status(200).JSON(
		fiber.Map{
			"status":      "success",
			"event_start": time.Unix(startTime, 0),
			"event_end":   time.Unix(endTime, 0),
			"post_event":  postEvent,
			"ctf_name":    ctfName,
			"team_len":    teamLen,
		},
	)
}
