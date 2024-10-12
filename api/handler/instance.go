package handler

import (
	"strconv"

	"github.com/TheAlpha16/isolet/api/deployment"
	"github.com/TheAlpha16/isolet/api/models"

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

	packed := new(models.AccessDetails)

	if err := deployment.DeployInstance(c, chall_id, teamid, packed); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "failure", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": packed})
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

// func GetStatus(c *fiber.Ctx) error {
// 	// var userid int
// 	// var err error

// 	// claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
// 	// userid = int(claims["userid"].(float64))

// 	// instances, err := database.GetInstances(c, userid)
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "contact admin"})
// 	// }
// 	return c.SendStatus(fiber.StatusOK)
// }
