package listener

import (
	"context"

	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	"github.com/TheAlpha16/isolet/oracle/utils/logger"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
)

func handleError(ctx context.Context, span trace.Span, log *logger.StandardLogger, msg string, err error, extraFields map[string]any) {
	span.RecordError(err)
	span.SetStatus(codes.Error, msg)
	errorDom.RaiseToSentry(ctx, err)

	zapFields := []zap.Field{zap.Error(err)}
	for k, v := range extraFields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	log.Error(msg, zapFields...)
}

func toInstance(ctx context.Context, obj interface{}) (*tidev1.Instance, bool, error) {
	if inst, ok := obj.(*tidev1.Instance); ok {
		return inst, false, nil
	}

	if tombstone, ok := obj.(cache.DeletedFinalStateUnknown); ok {
		if inst, ok := tombstone.Obj.(*tidev1.Instance); ok {
			return inst, true, nil
		}
		return nil, true, errorDom.RaiseInternal(ctx, "tombstone contained object that is not an Instance", nil, nil)
	}

	return nil, false, errorDom.RaiseInternal(ctx, "couldn't get object from event", nil, nil)
}
