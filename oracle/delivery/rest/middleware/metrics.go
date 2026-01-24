package middleware

import (
	"strconv"
	"time"

	"github.com/TheAlpha16/isolet/oracle/utils"
	"github.com/TheAlpha16/isolet/oracle/utils/tracer"

	"github.com/gofiber/fiber/v2"
	fiberutils "github.com/gofiber/fiber/v2/utils"
)

var skipMetricRoutes = map[string]struct{}{
	utils.RouteMetrics: {},
	utils.RoutePing:    {},
}

func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := fiberutils.CopyString(c.Method())
		path := fiberutils.CopyString(c.Path())

		if _, skip := skipMetricRoutes[path]; skip {
			return c.Next()
		}

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
