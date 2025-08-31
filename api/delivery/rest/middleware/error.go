package middleware

import (
	"fmt"

	"github.com/TheAlpha16/isolet/api/delivery/rest/response"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ErrorMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defaultHandler := c.App().Config().ErrorHandler

		defer func() {
			if r := recover(); r != nil {
				panicErr := errorDom.Raise(
					c.UserContext(), errorDom.ErrInternalServerError,
					"", fmt.Errorf("%v", r),
					errorDom.ExtraData{"stacktrace": errorDom.GetStackTrace()},
				)
				// log the error and send it to sentry
				logger.GetLogger(c.UserContext()).Error("recovered panic", zap.Error(panicErr))
				errorDom.RaiseToSentry(c.UserContext(), panicErr)

				// return a generic internal server error
				err = errorDom.Raise(
					c.UserContext(), errorDom.ErrInternalServerError,
					"", nil, nil,
				)
				_ = defaultHandler(c, err)
			}
		}()
		err = c.Next()
		if err != nil {
			code := errorDom.GetHTTPStatusCode(err)
			if code == -1 {
				// let fiber handle unrecognized errors
				return defaultHandler(c, err)
			}

			// filter out internal errors
			message := err.Error()
			if code == fiber.StatusInternalServerError {
				message = "Internal Server Error"
				logger.GetLogger(c.UserContext()).Error(message, zap.Error(err))
				errorDom.RaiseToSentry(c.UserContext(), err)
			}

			return c.Status(code).JSON(
				response.Error[any](message, nil),
			)
		}
		return nil
	}
}
