package tracer

import (
	"context"

	"github.com/TheAlpha16/isolet/oracle/utils/logger"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func StartSpan(ctx context.Context, tracer trace.Tracer, spanName string) (context.Context, trace.Span, *zap.Logger) {
	ctx, span := tracer.Start(ctx, spanName)
	log := logger.GetAppLogger().With(
		zap.String("span", spanName),
		zap.String("span_id", span.SpanContext().SpanID().String()),
		zap.String("trace_id", span.SpanContext().TraceID().String()),
	)
	return ctx, span, log
}
