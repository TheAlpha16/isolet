package handler

import (
	// "encoding/json"
	// "strconv"
	"fmt"
	"log"
	"strings"

	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/utils"

	// "github.com/TheAlpha16/isolet/api/deployment"
	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetStatus(c *fiber.Ctx) error {
	// var userid int
	// var err error

	// claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	// userid = int(claims["userid"].(float64))

	// instances, err := database.GetInstances(c, userid)
	// if err != nil {
	// 	log.Println(err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "contact admin"})
	// }
	return c.SendStatus(fiber.StatusOK)
}

func GetChalls(c *fiber.Ctx) error {
	challenges, err := database.ReadChallenges(c)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in reading challenges"})
	}
	return c.Status(fiber.StatusOK).JSON(challenges)
}

func CreateTeam(c *fiber.Ctx) error {
	var userid int
	var teamid int
	team := new(models.Team)

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userid = int(claims["userid"].(float64))
	teamid = int(claims["teamid"].(float64))

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

	if database.TeamExists(team.TeamName) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "failure", "message": "teamname already taken"})
	}

	team.Password = utils.Hash(team.Password)

	if err := database.CreateTeam(c, team); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in creating team, contact admin"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": fmt.Sprintf("team created successfully '%s' '%s' '%d'", team.TeamName, team.Password, teamid)})
}

func JoinTeam(c *fiber.Ctx) error {
	var userid int
	var teamid int
	team := new(models.Team)

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userid = int(claims["userid"].(float64))
	teamid = int(claims["teamid"].(float64))

	if teamid != -1 || database.UserInTeam(c, userid) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "failure", "message": "user already in a team"})
	}

	team.TeamName = strings.TrimSpace(c.FormValue("teamname"))
	team.Password = strings.TrimSpace(c.FormValue("password"))

	if team.TeamName == "" || team.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing teamname or password in request"})
	}

	if !database.TeamExists(team.TeamName) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "failure", "message": "team does not exist"})
	}

	team.Password = utils.Hash(team.Password)

	if err := database.AuthenticateTeam(c, team); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "invalid team credentials"})
	}

	if err := database.JoinTeam(c, team.TeamName, userid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "joined team successfully"})
}

// func StartInstance(c *fiber.Ctx) error {
// 	var userid int
// 	var err error

// 	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
// 	userid = int(claims["userid"].(float64))

// 	level_string := c.FormValue("level")

// 	if level_string == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing level in request"})
// 	}

// 	level, err := strconv.Atoi(level_string)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid level"})
// 	}

// 	if !database.UserExists(c, userid) {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "user does not exist"})
// 	}

// 	if !database.ValidChallenge(c, level) {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "level does not exist"})
// 	}

// 	if !database.CanStartInstance(c, userid, level) {
// 		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "failure", "message": "concurrent instances limit reached"})
// 	}

// 	deadline, password, port, hostname, err := deployment.DeployInstance(c, userid, level)
// 	if err != nil {
// 		database.DeleteRunning(c, userid, level)
// 		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "failure", "message": "error in initiating instance, contact admin"})
// 	}

// 	packed, err := json.Marshal(models.AccessDetails{Password: password, Port: port, Hostname: hostname, Deadline: deadline})
// 	if err != nil {
// 		log.Println(err)
// 		deployment.DeleteInstance(c, userid, level)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in initiating instance, contact admin"})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": packed})
// }

// func StopInstance(c *fiber.Ctx) error {
// 	var userid int
// 	var err error

// 	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
// 	userid = int(claims["userid"].(float64))

// 	level_string := c.FormValue("level")

// 	if level_string == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing parameters in request"})
// 	}

// 	level, err := strconv.Atoi(level_string)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid level"})
// 	}

// 	if !database.UserExists(c, userid) {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "user does not exist"})
// 	}

// 	if !database.ValidFlagEntry(c, level, userid) {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "instance stopped, reload page"})
// 	}

// 	if err := deployment.DeleteInstance(c, userid, level); err != nil {
// 		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "failure", "message": "error in stopping instance, contact admin"})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "instance stopped successfully"})
// }

// func SubmitFlag(c *fiber.Ctx) error {
// 	var userid int
// 	var err error

// 	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
// 	userid = int(claims["userid"].(float64))

// 	level_string := c.FormValue("level")
// 	flag := c.FormValue("flag")

// 	if flag == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing flag in the request"})
// 	}

// 	if level_string == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing parameters in request"})
// 	}

// 	level, err := strconv.Atoi(level_string)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid level"})
// 	}
// 	flag = strings.TrimSpace(flag)

// 	if !database.UserExists(c, userid) {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "user does not exist"})
// 	}

// 	if !database.ValidFlagEntry(c, level, userid) {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "instance not running"})
// 	}

// 	if isOK, message := database.VerifyFlag(c, level, userid, flag); !isOK {
// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "failure", "message": message})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "correct flag"})
// }

// func ExtendTime(c *fiber.Ctx) error {
// 	var userid int
// 	var err error

// 	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
// 	userid = int(claims["userid"].(float64))

// 	level_string := c.FormValue("level")

// 	if level_string == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing level in request"})
// 	}

// 	level, err := strconv.Atoi(level_string)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid level"})
// 	}

// 	if !database.UserExists(c, userid) {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "user does not exist"})
// 	}

// 	if !database.ValidFlagEntry(c, level, userid) {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "instance not running"})
// 	}

// 	isOK, message, newdeadline := deployment.AddTime(c, userid, level)
// 	if !isOK {
// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "failure", "message": message})
// 	}

// 	packed, err := json.Marshal(models.ExtendDeadline{Deadline: newdeadline})
// 	if err != nil {
// 		log.Println(err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in extension, contact admin"})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": packed})
// }

// func ShowScoreBoard(c *fiber.Ctx) error {
// 	board, err := database.ReadScores(c)
// 	if err != nil {
// 		log.Println(err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in reading scores"})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(board)
// }
