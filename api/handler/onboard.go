package handler

import (
	"log"
	"time"
	"strings"

	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/middleware"
	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)


func CreateTeam(c *fiber.Ctx) error {
	var userid int64
	var teamid int
	var email string
	var rank int
	team := new(models.Team)
	user := new(models.User)

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userid = int64(claims["userid"].(float64))
	email = claims["email"].(string)
	teamid = int(claims["teamid"].(float64))
	rank = int(claims["rank"].(float64))

	if teamid != -1 || database.UserInTeam(c, userid) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "failure", "message": "user already in a team"})
	}

	team.TeamName = strings.TrimSpace(c.FormValue("teamname"))
	team.Password = strings.TrimSpace(c.FormValue("password"))
	team.Captain = userid

	if team.TeamName == "" || team.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing teamname or password in request"})
	}

	if len(team.TeamName) < 3 || len(team.TeamName) > 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "teamname should be between 3 and 32 characters"})
	}

	if len(team.Password) < 6 || len(team.Password) > 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "password should be between 8 and 32 characters"})
	}

	if database.TeamNameExists(team.TeamName) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "failure", "message": "teamname already taken"})
	}

	team.Password = utils.Hash(team.Password)

	if err := database.CreateTeam(c, team); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in creating team, contact admin"})
	}

	user.UserID = int64(userid)
	user.Email = email
	user.Rank = rank
	user.TeamID = team.TeamID

	token, err := middleware.GenerateToken(user)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in token generation. contact admin"})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.SameSite = fiber.CookieSameSiteStrictMode
	cookie.Expires = time.Now().Add(72 * time.Hour)
	c.Cookie(cookie)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"teamid": team.TeamID,
	})
}

func JoinTeam(c *fiber.Ctx) error {
	var userid int64
	var teamid int64
	var email string
	var rank int
	team := new(models.Team)
	user := new(models.User)

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userid = int64(claims["userid"].(float64))
	email = claims["email"].(string)
	teamid = int64(claims["teamid"].(float64))
	rank = int(claims["rank"].(float64))

	if teamid != -1 || database.UserInTeam(c, userid) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "failure", "message": "user already in a team"})
	}

	team.TeamName = strings.TrimSpace(c.FormValue("teamname"))
	team.Password = strings.TrimSpace(c.FormValue("password"))

	if team.TeamName == "" || team.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing teamname or password in request"})
	}

	if !database.TeamNameExists(team.TeamName) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "failure", "message": "team does not exist"})
	}

	team.Password = utils.Hash(team.Password)

	if err := database.AuthenticateTeam(c, team); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "invalid team credentials"})
	}

	if err := database.JoinTeam(c, team.TeamName, userid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": err.Error()})
	}

	user.UserID = int64(userid)
	user.Email = email
	user.Rank = rank
	user.TeamID = team.TeamID

	token, err := middleware.GenerateToken(user)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in token generation. contact admin"})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.SameSite = fiber.CookieSameSiteStrictMode
	cookie.Expires = time.Now().Add(72 * time.Hour)
	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"teamid": team.TeamID,
	})
}
