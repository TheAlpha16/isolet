package errors

import (
	"context"

	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel/trace"
)

// ExtractErrorCode extracts the ErrorCode from an error, if it is an AppError
func ExtractErrorCode(err error) (ErrorCode, bool) {
	if ae, ok := AsAppError(err); ok {
		return ae.ErrorCode, true
	}

	return "", false
}

// IsSameError checks if the error is the same as the given ErrorCode
func IsSameError(err error, code ErrorCode) bool {
	if ec, ok := ExtractErrorCode(err); ok {
		return ec == code
	}
	return false
}

func addTraceContextToSentryEvent(ctx context.Context, event *sentry.Event) {
	span := trace.SpanFromContext(ctx)
	if span == nil || !span.SpanContext().IsValid() {
		return
	}
	sc := span.SpanContext()
	event.Contexts["trace"] = map[string]interface{}{
		"trace_id": sc.TraceID().String(),
		"span_id":  sc.SpanID().String(),
		"type":     "trace",
	}
}
