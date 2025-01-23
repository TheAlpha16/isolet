package handler

import (
	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {
	return c.SendString("Profile")
}