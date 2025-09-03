package errors

import (
	"context"
	"errors"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"

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

// IsServerSideError checks if the error is a server-side error
func IsServerSideError(err error) bool {
	errCode, ok := ExtractErrorCode(err)
	if !ok {
		return false
	}
	_, isServerSide := ServerSideErrors[errCode]
	return isServerSide
}

// IsAuthError checks if the error is an auth error
func IsAuthError(err error) bool {
	errCode, ok := ExtractErrorCode(err)
	if !ok {
		return false
	}
	_, isAuth := AuthErrors[errCode]
	return isAuth
}

// IsForbiddenError checks if the error is a forbidden error
func IsForbiddenError(err error) bool {
	errCode, ok := ExtractErrorCode(err)
	if !ok {
		return false
	}
	_, isForbidden := ForbiddenErrors[errCode]
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
