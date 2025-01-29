package handler

import (
	"fmt"

	"github.com/TheAlpha16/isolet/admin/database"
	"github.com/TheAlpha16/isolet/admin/models"
	"github.com/TheAlpha16/isolet/admin/utils"
	"github.com/gofiber/fiber/v2"
)

func EditChallengeMetaData(c *fiber.Ctx) error {
	var challengeMetadata models.Challenge
	if err := c.BodyParser(&challengeMetadata); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failure",
			"message": "invalid request body",
		})
	}

	if err := utils.ValidateChallengeFields(&challengeMetadata); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failure",
			"message": err.Error(),
		})
	}

	// fetch existing challenge
	existingChallenge, err := database.FetchChallenge(c, challengeMetadata.ChallID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failure",
			"message": err.Error(),
		})
	}

	fmt.Printf("existing challenge: %v\n", existingChallenge)

	// update challenge field properties
	updatedChallenge:= database.UpdateChallenges(c, &existingChallenge, &challengeMetadata)

	// save challenge
	if err := database.SaveChallengeMetaData(c, updatedChallenge); err != nil {
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

func EditChallengeHints(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}

func EditChallengeFiles(c *fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}

