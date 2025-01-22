package adminhandler

import "github.com/gofiber/fiber/v2"

func EditChallenges(c *fiber.Ctx) error {



	return c.SendStatus(fiber.StatusOK)
}