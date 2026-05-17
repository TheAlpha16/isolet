package errors

import (
	"context"

	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"

	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ExtractErrorCode extracts the ErrorCode from an error, if it is an AppError
func ExtractErrorCode(err error) (ErrorCode, bool) {
	if ae, ok := AsAppError(err); ok {
		return ae.ErrorCode, true
	}

	return "", false
}

// Is checks if the error is the same as the given ErrorCode
func Is(err error, code ErrorCode) bool {
	if ec, ok := ExtractErrorCode(err); ok {
		return ec == code
	}
	return false
}

// IsServerErrorCode checks if the error is a server-side error
func IsServerErrorCode(code ErrorCode) bool {
	_, isServerSide := ServerErrorCodes[code]
	return isServerSide
}

// IsAuthErrorCode checks if the error is an auth error
func IsAuthErrorCode(code ErrorCode) bool {
	_, isAuth := AuthErrorCodes[code]
	return isAuth
}

// IsForbiddenErrorCode checks if the error is a forbidden error
func IsForbiddenErrorCode(code ErrorCode) bool {
	_, isForbidden := ForbiddenErrorCodes[code]
	return isForbidden
}

// Sets `ExtraData` from context in the app error
func EnrichWithCtx(ctx context.Context, ae *AppError) {
	extraData := common.GetExtraDataFromCtx(ctx)
	for k, v := range extraData {
		ae.AddExtraData(k, v)
	}
}

// Sets `Src` in the given context.
func SetSrcInCtx(ctx context.Context, src string) context.Context {
	return context.WithValue(ctx, srcKey, src)
}

// Extracts `Src` from the given context
func GetSrcFromCtx(ctx context.Context) string {
	ctxVal := ctx.Value(srcKey)
	switch src := ctxVal.(type) {
	case string:
		return src
	default:
		return ""
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

func HandleSpanError(ctx context.Context, span trace.Span, log *zap.Logger, msg string, err error, extraFields ...zap.Field) {
	span.RecordError(err)
	span.SetStatus(codes.Error, msg)
	RaiseToSentry(ctx, err)

	zapFields := append(make([]zap.Field, 0, 1+len(extraFields)), zap.Error(err))
	zapFields = append(zapFields, extraFields...)
	log.Error(msg, zapFields...)
}
