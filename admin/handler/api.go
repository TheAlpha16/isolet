package handler

import (
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
			"message": err.Error(),
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

	// update challenge field properties
	updatedChallenge:= utils.UpdateChallenge(&existingChallenge, &challengeMetadata)

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

func EditChallengeFiles(c *fiber.Ctx) error {
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

	// update file field properties
	updatedChallenge:= utils.UpdateFiles(&existingChallenge, &challengeMetadata)

	// save challenge
	if err := database.SaveChallengeMetaData(c, updatedChallenge); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failure",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"message": "challenge files updated successfully",
	})
}

func EditChallengeHints(c *fiber.Ctx) error {
	var hintData models.Hint
	if err := c.BodyParser(&hintData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "failure",
			"message": "invalid request body",
		})
	}

	if err := utils.ValidateHintFields(&hintData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "failure",
			"message": err.Error(),
		})
	}

	// fetch existing hints
	existingHint, err := database.FetchHint(c, hintData.HID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "failure",
			"message": err.Error(),
		})
	}

	// update hints data
	updatedHint := utils.UpdateHint(&existingHint, &hintData)
	
	// save hints
	if err := database.SaveHintData(c, updatedHint); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "failure",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"message": "challenge hints updated successfully",
	})
}

func EditChallengeRequirements(c *fiber.Ctx) error {
	var challengeMetadata models.Challenge
	if err := c.BodyParser(&challengeMetadata); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failure",
			"message": "invalid request body",
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

	// update file field properties
	updatedChallenge:= utils.UpdateRequirements(&existingChallenge, &challengeMetadata)


	// save challenge
	if err := database.SaveChallengeMetaData(c, updatedChallenge); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failure",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"message": "challenge requirements updated successfully",
	})
}

