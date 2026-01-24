package middleware

import (
	"strconv"
	"time"

	"github.com/TheAlpha16/isolet/oracle/utils/tracer"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := utils.CopyString(c.Method())
		path := utils.CopyString(c.Path())

		tracer.HttpRequestsInFlight.WithLabelValues(path).Inc()
		defer tracer.HttpRequestsInFlight.WithLabelValues(path).Dec()

		err := c.Next()

		tracer.HttpRequestsTotal.WithLabelValues(
			method,
			path,
			strconv.Itoa(c.Response().StatusCode()),
		).Inc()
		tracer.HttpRequestDurationSeconds.WithLabelValues(
			method,
			path,
		).Observe(time.Since(start).Seconds())

		return err
	}
}
