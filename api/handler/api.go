package handler

import (
	"log"
	"strconv"
	"strings"

	"github.com/TheAlpha16/isolet/api/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)


func GetChalls(c *fiber.Ctx) error {
	var teamid int64

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	teamid = int64(claims["teamid"].(float64))

	challenges, err := database.ReadChallenges(c, teamid)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in reading challenges"})
	}
	return c.Status(fiber.StatusOK).JSON(challenges)
}

func SubmitFlag(c *fiber.Ctx) error {
	var userid int64
	var teamid int64
	var err error

	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userid = int64(claims["userid"].(float64))
	teamid = int64(claims["teamid"].(float64))

	chall_id_string := c.FormValue("chall_id")
	flag := c.FormValue("flag")

	if flag == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing flag in the request"})
	}

	if chall_id_string == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "missing chall_id in the request"})
	}

	chall_id, err := strconv.Atoi(chall_id_string)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "invalid chall_id"})
	}

	flag = strings.TrimSpace(flag)

	if isOK, message := database.VerifyFlag(c, chall_id, userid, teamid, flag); !isOK {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "failure", "message": message})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "correct flag"})
}

// func ShowScoreBoard(c *fiber.Ctx) error {
// 	board, err := database.ReadScores(c)
// 	if err != nil {
// 		log.Println(err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in reading scores"})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(board)
// }
