package handler

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/admin/database"
	"github.com/TheAlpha16/isolet/admin/middleware"
	"github.com/TheAlpha16/isolet/admin/models"
	"github.com/TheAlpha16/isolet/admin/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

func Register(c *fiber.Ctx) error {
	regForm := new(models.ToVerify)

	regForm.Email = c.FormValue("email")
	regForm.Username = c.FormValue("username")
	regForm.Password = c.FormValue("password")
	regForm.Confirm = c.FormValue("confirm")
	if regForm.Email == "" || regForm.Username == "" || regForm.Password == "" || regForm.Confirm == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "all fields are required"})
	}

	regForm.Email = strings.TrimSpace(regForm.Email)
	regForm.Username = strings.TrimSpace(regForm.Username)
	regForm.Password = strings.TrimSpace(regForm.Password)
	regForm.Confirm = strings.TrimSpace(regForm.Confirm)

	regForm.Email = strings.ToLower(regForm.Email)
	regForm.Username = strings.ToLower(regForm.Username)

	if isOK, status := utils.ValidateRegisterInput(regForm); !isOK {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": status})
	}

	if err := utils.SendVerificationMail(regForm); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in sending verification mail"})
	}

	regForm.Password = utils.Hash(regForm.Password)

	if err := database.AddToVerify(c, regForm); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "please contact admin"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "check your mail for verification"})
}

func Verify(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		c.Status(fiber.StatusBadRequest).SendString("missing token. Register again!")
	}
	claims := new(models.VerifyClaims)

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.TOKEN_SECRET), nil
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if claims.Email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("error in token, register again!")
	}

	if message, err := database.AddToUsers(c, claims.Email); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString(message)
	}

	return c.Status(fiber.StatusCreated).SendString("user verified successfully! proceed to login")
}
