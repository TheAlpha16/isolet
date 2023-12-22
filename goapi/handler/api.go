package handler

import (
	// "fmt"
	"encoding/json"
	"log"
	"strconv"

	"github.com/TitanCrew/isolet/config"
	"github.com/TitanCrew/isolet/database"
	"github.com/TitanCrew/isolet/deployment"
	"github.com/TitanCrew/isolet/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetStatus(c *fiber.Ctx) error {
	var userid int
	var err error

	if !config.DISCORD_FRONTEND {
		claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
		userid = int(claims["userid"].(float64))
	} else {
		userid_string := c.FormValue("userid")
		if userid_string == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing parameters in request"})
		}
		userid, err = strconv.Atoi(userid_string)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid userid"})
		}
	}

	instances, err := database.GetInstances(userid)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "contact admin"})
	}
	return c.Status(fiber.StatusOK).JSON(instances)
}

func GetChalls(c *fiber.Ctx) error {
	challenges, err := database.ReadChallenges()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in reading challenges"})
	}
	return c.Status(fiber.StatusOK).JSON(challenges)
}

func StartInstance(c *fiber.Ctx) error {
	var userid int
	var err error

	if !config.DISCORD_FRONTEND {
		claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
		userid = int(claims["userid"].(float64))
	} else {
		userid_string := c.FormValue("userid")
		if userid_string == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing parameters in request"})
		}
		userid, err = strconv.Atoi(userid_string)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid userid"})
		}
	}

	level_string := c.FormValue("level")
	chall_id_string := c.FormValue("chall_id")

	if level_string == "" || chall_id_string == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing parameters in request"})
	}

	level, err := strconv.Atoi(level_string)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid level"})
	}

	chall_id, err := strconv.Atoi(chall_id_string)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid chall_id"})
	}

	if !database.UserExists(userid) {
		if !config.DISCORD_FRONTEND {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "user does not exist"})
		}
		database.AddToUsersDiscord(userid)
	}

	if !database.ValidChallenge(chall_id, level) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid level"})
	}

	if !database.CanStartInstance(userid, level) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "failure", "message": "concurrent instances limit reached"})
	}

	password, port, err := deployment.DeployInstance(userid, level)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "failure", "message": "error in initiating instance, contact admin"})
	}

	packed, err := json.Marshal(models.AccessDetails{Password: password, Port: port})
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in initiating instance, contact admin"})
		
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": packed})
}

func StopInstance(c *fiber.Ctx) error {
	var userid int
	var err error

	if !config.DISCORD_FRONTEND {
		claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
		userid = int(claims["userid"].(float64))
	} else {
		userid_string := c.FormValue("userid")
		if userid_string == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing parameters in request"})
		}
		userid, err = strconv.Atoi(userid_string)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid userid"})
		}
	}

	level_string := c.FormValue("level")

	if level_string == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing parameters in request"})
	}

	level, err := strconv.Atoi(level_string)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid level"})
	}

	if !database.UserExists(userid) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "user does not exist"})
	}

	if !database.ValidFlagEntry(level, userid) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "instance stopped, reload page"})
	}

	if err := deployment.DeleteInstance(userid, level); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "failure", "message": "error in initiating instance, contact admin"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "instance stopped successfully"})
}

func SubmitFlag(c *fiber.Ctx) error {
	var userid int
	var err error

	if !config.DISCORD_FRONTEND {
		claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
		userid = int(claims["userid"].(float64))
	} else {
		userid_string := c.FormValue("userid")
		if userid_string == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing parameters in request"})
		}
		userid, err = strconv.Atoi(userid_string)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid userid"})
		}
	}

	level_string := c.FormValue("level")
	flag := c.FormValue("flag")

	if flag == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing flag in the request"})
	}

	if level_string == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing parameters in request"})
	}

	level, err := strconv.Atoi(level_string)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid level"})
	}

	if !database.UserExists(userid) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "user does not exist"})
	}

	if !database.ValidFlagEntry(level, userid) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "instance not running"})
	}

	if isOK, message := database.VerifyFlag(level, userid, flag); !isOK {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "failure", "message": message})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "correct flag"})
}

func ShowScoreBoard(c *fiber.Ctx) error {
	board, err := database.ReadScores()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in reading scores"})
	}
	return c.Status(fiber.StatusOK).JSON(board)
}