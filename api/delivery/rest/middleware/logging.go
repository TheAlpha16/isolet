package middleware

import (
	"time"

	"github.com/TheAlpha16/isolet/api/utils/logger"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var skipPaths = map[string]struct{}{
	"/ping": {},
}

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if _, ok := skipPaths[c.Path()]; ok {
			return c.Next()
		}

		start := time.Now()
		ctx := c.UserContext()
		span := trace.SpanFromContext(ctx)
		childLogger := logger.GetLogger(ctx).WithFields(
			zap.String("span_id", span.SpanContext().SpanID().String()),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)
		ctx = logger.NewContext(ctx, childLogger)
		c.SetUserContext(ctx)

		childLogger.Info("received request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
		)

		err := c.Next()
		if err != nil {
			return err
		}

		childLogger.Info("request completed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("latency", time.Since(start).Truncate(time.Millisecond).String()),
		)

		return nil
	}
}
