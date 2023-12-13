package middleware

import (
	"time"

	"github.com/TitanCrew/isolet/config"
	"github.com/TitanCrew/isolet/models"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CheckToken() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(config.SESSION_SECRET),
			JWTAlg: jwtware.HS256,
		},
	})
}

func GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"userid": user.UserID,
		"email": user.Email,
		"rank": user.Rank,
		"exp": time.Now().Add(time.Hour * time.Duration(config.SESSION_EXP)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.SESSION_SECRET))
	if err != nil {
		return "", err
	}
	return t, nil
}