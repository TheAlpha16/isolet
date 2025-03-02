package router

import (
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

var RateLimits = map[string]limiter.Config{
	AUTH_GROUP + METADATA:        NewLimitConfig(10, 10*time.Second),
	AUTH_GROUP + LOGIN:           NewLimitConfig(5, 10*time.Second),
	AUTH_GROUP + REGISTER:        NewLimitConfig(5, 10*time.Second),
	AUTH_GROUP + VERIFY:          NewLimitConfig(5, 10*time.Second),
	AUTH_GROUP + FORGOT_PASSWORD: NewLimitConfig(5, 10*time.Second),
	AUTH_GROUP + RESET_PASSWORD:  NewLimitConfig(5, 10*time.Second),

	ONBOARD_GROUP + TEAM_JOIN:        NewLimitConfig(3, 1*time.Minute),
	ONBOARD_GROUP + TEAM_CREATE:      NewLimitConfig(5, 30*time.Second),
	ONBOARD_GROUP + JOIN_TEAM_INVITE: NewLimitConfig(5, 30*time.Second),

	API_GROUP + CHALLS:      NewLimitConfig(10, 10*time.Second),
	API_GROUP + SUBMIT:      NewLimitConfig(5, 10*time.Second),
	API_GROUP + UNLOCK_HINT: NewLimitConfig(5, 10*time.Second),

	API_GROUP + SCOREBOARD_GROUP: NewLimitConfig(7, 10*time.Second),

	API_GROUP + IDENTIFY: NewLimitConfig(7, 10*time.Second),
	API_GROUP + LOGOUT:   NewLimitConfig(7, 10*time.Second),

	API_GROUP + LAUNCH: NewLimitConfig(5, 10*time.Second),
	API_GROUP + STOP:   NewLimitConfig(5, 10*time.Second),
	API_GROUP + STATUS: NewLimitConfig(5, 10*time.Second),
	API_GROUP + EXTEND: NewLimitConfig(5, 10*time.Second),

	API_GROUP + PROFILE_GROUP: NewLimitConfig(15, 10*time.Second),
}

func RateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		for route, config := range RateLimits {
			if strings.HasPrefix(c.Path(), route) {
				config.Storage = utils.GetActiveStore()
				return limiter.New(config)(c)
			}
		}

		return c.Next()
	}
}

func NewLimitConfig(max int, expiration time.Duration) limiter.Config {
	return limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			rlEnabled, err := config.GetBool("API_RATE_LIMIT")
			if err != nil {
				return false
			}
			return !rlEnabled
		},
		Max:        max,
		Expiration: expiration,
		Storage:    utils.GetActiveStore(),
		KeyGenerator: func(c *fiber.Ctx) string {
			var TeamNameKey models.TeamNameKey
			if team, ok := c.Locals(TeamNameKey).(string); ok {
				return "rate:team:" + team
			}
			return "rate:ip:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "failure",
				"message": "rate limit exceeded",
			})
		},
		LimiterMiddleware: limiter.SlidingWindow{},
	}
}
