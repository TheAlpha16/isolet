package middleware

import (
	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	"github.com/gofiber/fiber/v2"
)

func CheckTimingsMiddleware(cvUc cvDom.Usecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
