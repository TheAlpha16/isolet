package errors

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	goerrors "github.com/go-errors/errors"
	"github.com/gofiber/fiber/v2"
)

// ErrorCode type
type ErrorCode string

// ErrCtxKey is a key for storing/retrieving error context data from context.Context
type ErrCtxKey string

// ExtraData type
type ExtraData map[string]any

// AppError represents an application error with a code, message, cause, and extra data
type AppError struct {
	ErrorCode ErrorCode // unique error code
	Message   string    // optional human-readable message
	Cause     error     // wrapped error
	ExtraData ExtraData // extra data associated with the error
}

func (ae *AppError) Error() string {
	baseErr := fmt.Sprintf("[%s] %s", ae.GetCode(), ae.GetMessage())
	if ae.Cause != nil {
		baseErr = fmt.Sprintf("%s : %s", baseErr, ae.Cause.Error())
	}
	return baseErr
}

func (ae *AppError) Unwrap() error {
	return ae.Cause
}

func (ae *AppError) GetErrorCode() ErrorCode {
	return ae.ErrorCode
}

func (ae *AppError) GetCode() string {
	return fmt.Sprintf("%v-%v", SVC, ae.ErrorCode)
}

func (ae *AppError) GetMessage() string {
	return ae.Message
}

func (ae *AppError) AddExtraData(key string, val any) {
	if ae.ExtraData == nil {
		ae.ExtraData = make(map[string]any)
	}
	ae.ExtraData[key] = val
}

func (ae *AppError) GetExtraData() ExtraData {
	return ae.ExtraData
}

func Raise(ctx context.Context, code ErrorCode, msg string, errToWrap error, extraData ExtraData) error {
	ae := &AppError{
		ErrorCode: code,
		ExtraData: extraData,
	}
	if msg != "" {
		ae.Message = msg
	} else {
		ae.Message = msgMap[code]
	}
	if errToWrap != nil {
		ae.Cause = errToWrap
	}
	if ctx != nil {
		EnrichWithCtx(ctx, ae)
	}
	return ae
}

// RaiseToSentry raises the error to Sentry
func RaiseToSentry(ctx context.Context, err error) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		sentry.CaptureException(err)
		return
	}
	event := sentry.NewEvent()
	event.Level = sentry.LevelError
	switch e := err.(type) {
	case *AppError:
		hub.WithScope(func(scope *sentry.Scope) {
			EnrichWithCtx(ctx, e)
			event.Message = fmt.Sprintf("%s in %s: %s", e.GetCode(), GetSrcFromCtx(ctx), e.GetMessage())
			event.Exception = []sentry.Exception{{
				Value:      fmt.Sprintf("%s: %s", e.GetCode(), GetSrcFromCtx(ctx)),
				Type:       fmt.Sprintf("%T", e),
				Stacktrace: sentry.ExtractStacktrace(e.Cause),
			}}
			event.Extra = map[string]interface{}{
				"context": e.GetExtraData(),
			}
		})
		hub.CaptureEvent(event)
	default:
		stackErr := goerrors.Wrap(e, 1)
		hub.WithScope(func(scope *sentry.Scope) {
			event.Message = fmt.Sprintf("%s in %s", e.Error(), GetSrcFromCtx(ctx))
			event.Exception = []sentry.Exception{{
				Value:      e.Error(),
				Type:       fmt.Sprintf("%T", e),
				Stacktrace: sentry.ExtractStacktrace(stackErr),
			}}
			event.Extra = map[string]interface{}{
				"context": GetExtraDataFromCtx(ctx),
			}
		})
		addTraceContextToSentryEvent(ctx, event)
		hub.CaptureEvent(event)
	}
}

func GetHTTPStatusCode(err error) int {
	_, isDomErr := ExtractErrorCode(err)
	if !isDomErr {
		return -1 // return -1 if not a domain error
	}
	if isServerSideErr := IsServerSideError(err); isServerSideErr {
		return fiber.StatusInternalServerError
	}
	if isAuthErr := IsAuthError(err); isAuthErr {
		return fiber.StatusUnauthorized
	}
	if isForbiddenErr := IsForbiddenError(err); isForbiddenErr {
		return fiber.StatusForbidden
	}
	return fiber.StatusBadRequest
}
