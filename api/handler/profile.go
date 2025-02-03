package handler

import (
	"fmt"

	"github.com/TheAlpha16/isolet/api/config"
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

	rank, score, err := database.GetTeamRank(c, teamid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in reading team"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"team":        team,
		"submissions": submissions,
		"rank":        rank,
		"score":       score,
	})
}

func GetInviteToken(c *fiber.Ctx) error {
	var userid int64

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userid = int64(claims["userid"].(float64))

	token, err := database.GenerateInviteToken(c, userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"token":  fmt.Sprintf("http://%s/onboard/team/invite?token=%s", config.PUBLIC_URL, token),
	})
}
