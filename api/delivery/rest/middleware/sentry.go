package middleware

import (
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func SentryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
			c.SetUserContext(ctx)
		}

		// Try to pick the remote parent span context first
		span := trace.SpanFromContext(ctx)
		if spanCtx := span.SpanContext(); spanCtx.IsValid() {
			// If available, set the remote trace and span ID
			hub.Scope().SetContext("trace", sentry.Context{
				"trace_id": spanCtx.TraceID().String(),
				"span_id":  spanCtx.SpanID().String(),
				"op":       "rest.server",
				"sampled":  spanCtx.IsSampled(),
			})
		}

		err := c.Next()

		// Update span status
		if span != nil {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
			}
		}

		return err
	}
}
