package handler

import (
	"github.com/TheAlpha16/isolet/admin/database"
	"github.com/TheAlpha16/isolet/admin/models"
	"github.com/TheAlpha16/isolet/admin/utils"
	"github.com/gofiber/fiber/v2"
)

func EditChallengeMetaData(c *fiber.Ctx) error {
	var challenge models.Challenge
	if err := c.BodyParser(&challenge); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failure",
			"message": "invalid request body",
		})
	}

	if err := utils.ValidateChallengeFields(&challenge); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failure",
			"message": err.Error(),
		})
	}

	if err := database.EditChallengeData(c, &challenge); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failure",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"message": "challenge metadata updated successfully",
	})
}

func EditChallengesHints(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}

func EditChallengeFiles(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}

