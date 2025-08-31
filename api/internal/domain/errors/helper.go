package errors

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel/trace"
)

// ExtractErrorCode extracts the ErrorCode from an error, if it is an AppError
func ExtractErrorCode(err error) (ErrorCode, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
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

func ExtraDataFromCtx(ctx context.Context) ExtraData {
	ctxVal := ctx.Value(errExtraDataKey)
	switch ex := ctxVal.(type) {
	case ExtraData:
		return ex
	default:
		return make(ExtraData)
	}
}

func EnrichWithCtx(ctx context.Context, ae *AppError) {
	extraData := ExtraDataFromCtx(ctx)
	for k, v := range extraData {
		ae.AddExtraData(k, v)
	}
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
