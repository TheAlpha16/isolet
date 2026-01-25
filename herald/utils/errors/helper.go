package errors

import (
	"context"

	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

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

func HandleSpanError(ctx context.Context, span trace.Span, log *zap.Logger, msg string, err error, extraFields ...zap.Field) {
	span.RecordError(err)
	span.SetStatus(codes.Error, msg)
	RaiseToSentry(ctx, err)

	zapFields := []zap.Field{zap.Error(err)}
	zapFields = append(zapFields, extraFields...)
	log.Error(msg, zapFields...)
}
