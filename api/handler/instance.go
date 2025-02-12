package handler

import (
	"strconv"

	"github.com/TheAlpha16/isolet/api/models"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/deployment"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func StartInstance(c *fiber.Ctx) error {
	var teamid int64
	var err error

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	teamid = int64(claims["teamid"].(float64))

	chall_id_string := c.FormValue("chall_id")

	if chall_id_string == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing challenge id in the request"})
	}

	chall_id, err := strconv.Atoi(chall_id_string)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid challenge id"})
	}

	instance, err := deployment.DeployInstance(c, chall_id, teamid)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "failure", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": instance})
}

func StopInstance(c *fiber.Ctx) error {
	var teamid int64
	var err error

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	teamid = int64(claims["teamid"].(float64))

	chall_id_string := c.FormValue("chall_id")

	if chall_id_string == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing challenge id in the request"})
	}

	chall_id, err := strconv.Atoi(chall_id_string)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid challenge id"})
	}

	if err := deployment.DeleteInstance(c, chall_id, teamid); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "failure", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "instance stopped successfully"})
}

func ExtendTime(c *fiber.Ctx) error {
	var teamid int64
	var err error

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	teamid = int64(claims["teamid"].(float64))

	chall_id_string := c.FormValue("chall_id")

	if chall_id_string == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing challenge id in the request"})
	}

	chall_id, err := strconv.Atoi(chall_id_string)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid challenge id"})
	}

	packed := new(models.ExtendDeadline)

	if err := deployment.AddTime(c, chall_id, teamid, packed); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "failure", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": packed})
}

func GetStatus(c *fiber.Ctx) error {
	var teamid int64
	var err error

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	teamid = int64(claims["teamid"].(float64))

	instances, err := database.GetInstances(c, teamid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "contact admin"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": instances})
}
