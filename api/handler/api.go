package handler

import (
	"time"
	"strconv"

	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ShowScoreBoard(c *fiber.Ctx) error {
	page_string := c.Query("page", "1")

	page, err := strconv.Atoi(page_string)
	if err != nil {
		page = 1
	}

	if page < 1 {
		page = 1
	}

	board, err := database.ReadScores(c, page)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "failure", 
			"message": "cannot retrieve scoreboard at the moment",
		})
	}

	return c.Status(fiber.StatusOK).JSON(board)
}

func Identify(c *fiber.Ctx) error {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userid := int64(claims["userid"].(float64))
	email := claims["email"].(string)
	username := claims["username"].(string)
	rank := int(claims["rank"].(float64))
	teamid := int64(claims["teamid"].(float64))

	var TeamNameKey models.TeamNameKey

	teamname := c.Locals(TeamNameKey)
	if teamname == nil {
		teamname = ""
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"userid": userid,
			"email": email,
			"username": username,
			"rank": rank,
			"teamid": teamid,
			"teamname": teamname,
		})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
	})
	
	return c.SendStatus(fiber.StatusOK)
}