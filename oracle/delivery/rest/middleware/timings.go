package middleware

import (
	"time"

	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"

	"github.com/gofiber/fiber/v2"
)

func CheckTimingsMiddleware(cvUc cvDom.Usecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		currentTime := int(time.Now().Unix())

		eventStart := cvUc.GetInt(ctx, cvDom.EventStart)
		eventEnd := cvUc.GetInt(ctx, cvDom.EventEnd)
		postEvent := cvUc.GetBool(ctx, cvDom.EventPostMode)

		if postEvent {
			return c.Next()
		}

		if currentTime < eventStart {
			return errorDom.Raise(ctx, errorDom.ErrEventNotStarted, "", nil, nil)
		}

		if currentTime > eventEnd {
			return errorDom.Raise(ctx, errorDom.ErrEventEnded, "", nil, nil)
		}

		return c.Next()
	}
}
