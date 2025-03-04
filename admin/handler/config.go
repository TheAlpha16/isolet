package handler

import (
	"github.com/TheAlpha16/isolet/admin/database"
	"github.com/TheAlpha16/isolet/admin/models"
	"github.com/TheAlpha16/isolet/admin/utils"
	"github.com/gofiber/fiber/v2"
)

func EditConfigValues(c *fiber.Ctx) error{
	var configData models.Config
	if err := c.BodyParser(&configData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failure",
			"message": err.Error(),
		})
	}

	if err := utils.ValidateConfigFields(&configData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failure",
			"message": err.Error(),
		})
	}

	// fetch existing config
	existingConfig, err := database.FetchConfig(c, configData.Key)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failure",
			"message": err.Error(),
		})
	}

	// update config field properties
	updatedChallenge:= utils.UpdateConfig(&existingConfig, &configData)

	// save config
	if err := database.SaveConfigData(c, updatedChallenge); err != nil {
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