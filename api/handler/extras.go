package handler

import (
	"github.com/TheAlpha16/isolet/api/config"
	"github.com/gofiber/fiber/v2"
)

func Ping(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Server is up!"})
}

func GetMetadata(c *fiber.Ctx) error {
	return c.Status(200).JSON(
		fiber.Map{
			"status":      "success",
			"event_start": config.EVENT_START,
			"event_end":   config.EVENT_END,
			"post_event":  config.POST_EVENT,
			"ctf_name":    config.CTF_NAME,
			"team_len":    config.TEAM_LEN,
		},
	)
}
