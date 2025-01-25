package handler

import "github.com/gofiber/fiber/v2"

func EditChallengeMetaData(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}

func EditChallengeFiles(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}

func EditChallengesHints(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}