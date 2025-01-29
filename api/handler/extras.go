package handler

import (
	"strconv"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/gofiber/fiber/v2"
)

func Ping(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Server is up!"})
}

func GetMetadata(c *fiber.Ctx) error {
	startTime, _ := strconv.ParseInt(config.EVENT_START, 10, 64)
	endTime, _ := strconv.ParseInt(config.EVENT_END, 10, 64)

	return c.Status(200).JSON(
		fiber.Map{
			"status":      "success",
			"event_start": time.Unix(startTime, 0),
			"event_end":   time.Unix(endTime, 0),
			"post_event":  config.POST_EVENT,
			"ctf_name":    config.CTF_NAME,
			"team_len":    config.TEAM_LEN,
		},
	)
}
