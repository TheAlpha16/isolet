package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ResolveIP() fiber.Handler {
    return func(c *fiber.Ctx) error {
        var ip string

        if cfIp := c.Get("CF-Connecting-IP"); cfIp != "" {
            ip = cfIp
        } else if xff := c.Get("X-Forwarded-For"); xff != "" {
            ip = strings.Split(xff, ",")[0]
        } else {
            ip = c.IP()
        }

        c.Locals("clientIP", ip)

        return c.Next()
    }
}