package handler

import (
	"errors"
	"log"
	"time"

	"github.com/CyberLabs-Infosec/isolet/goapi/config"
	"github.com/CyberLabs-Infosec/isolet/goapi/database"
	"github.com/CyberLabs-Infosec/isolet/goapi/middleware"
	"github.com/CyberLabs-Infosec/isolet/goapi/models"
	"github.com/CyberLabs-Infosec/isolet/goapi/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *fiber.Ctx) error {
	creds := new(models.Creds)
	user := new(models.User)

	creds.Email = c.FormValue("email")
	creds.Password = c.FormValue("password")
	if creds.Email == "" || creds.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "username and password required"})
	}

	isValid, message := utils.ValidateLoginInput(creds)
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": message})
	}

	creds.Password = utils.Hash(creds.Password)

	if err := database.ValidateCreds(creds, user); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failure", "message": "invalid credentials"})
	}

	token, err := middleware.GenerateToken(user)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in token generation. contact admin"})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.SameSite = fiber.CookieSameSiteStrictMode
	// cookie.HTTPOnly = true
	cookie.Expires = time.Now().Add(72 * time.Hour)
	// cookie.Secure = true change this
	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "login successful"})
}

func Register(c *fiber.Ctx) error {
	regForm := new(models.User)

	regForm.Email = c.FormValue("email")
	regForm.Username = c.FormValue("username")
	regForm.Password = c.FormValue("password")
	regForm.Confirm = c.FormValue("confirm")
	if regForm.Email == "" || regForm.Username == "" || regForm.Password == "" || regForm.Confirm == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": "all fields are required"})
	}

	if isOK, status := utils.ValidateRegisterInput(regForm); !isOK {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failure", "message": status})
	}

	if err := utils.SendVerificationMail(regForm); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failure", "message": "error in sending verification mail"})
	}

	regForm.Password = utils.Hash(regForm.Password)

	if err := database.AddToVerify(regForm); err != nil {
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

	if message, err := database.AddToUsers(claims.Email); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString(message)
	}

	return c.Status(fiber.StatusCreated).SendString("user verified successfully! proceed to login")
}