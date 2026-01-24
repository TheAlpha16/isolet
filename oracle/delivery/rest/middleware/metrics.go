package middleware

import (
	"fmt"
	"time"

	"github.com/TheAlpha16/isolet/oracle/utils/tracer"
	"github.com/gofiber/fiber/v2"
)

func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := c.Method()
		path := c.Route().Path

		tracer.HttpRequestsInFlight.WithLabelValues(path).Inc()
		defer tracer.HttpRequestsInFlight.WithLabelValues(path).Dec()

		err := c.Next()

		status := c.Response().StatusCode()
		duration := time.Since(start).Seconds()
		statusStr := fmt.Sprintf("%d", status)

		tracer.HttpRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
		tracer.HttpRequestDurationSeconds.WithLabelValues(method, path).Observe(duration)

		return err
	}
}
