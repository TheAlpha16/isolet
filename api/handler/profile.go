package handler

import (
	"github.com/TheAlpha16/isolet/api/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetSelfTeam(c *fiber.Ctx) error {
	var teamid int64

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	teamid = int64(claims["teamid"].(float64))

	team, err := database.ReadTeam(c, teamid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in reading team"})
	}

	submissions, err := database.GetSubmissions(c, teamid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in reading team"})
	}

	rank, err := database.GetTeamRank(c, teamid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in reading team"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"team":        team,
		"submissions": submissions,
		"rank":        rank,
	})
}
