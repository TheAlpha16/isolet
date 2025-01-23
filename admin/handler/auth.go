package handler

import (
	"log"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/admin/database"
	"github.com/TheAlpha16/isolet/admin/middleware"
	"github.com/TheAlpha16/isolet/admin/models"
	"github.com/TheAlpha16/isolet/admin/utils"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	user := new(models.User)

	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "username/email and password required"})
	}

	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)

	user.Email = strings.ToLower(user.Email)

	isValid, message := utils.ValidateLoginInput(user)
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": message})
	}

	user.Password = utils.Hash(user.Password)

	if err := database.ValidateCreds(c, user); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "invalid credentials"})
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in token generation. contact admin"})
	}

	teamname, _ := database.TeamExists(c, user.TeamID)

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.SameSite = fiber.CookieSameSiteStrictMode
	cookie.HTTPOnly = true
	cookie.Expires = time.Now().Add(72 * time.Hour)
	// cookie.Secure = true change this
	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"userid": user.UserID,
		"email": user.Email,
		"username": user.Username,
		"rank": user.Rank,
		"teamid": user.TeamID,
		"teamname": teamname,
	})
}
