package errors

import (
	"context"
	"fmt"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	goerrors "github.com/go-errors/errors"
	"github.com/gofiber/fiber/v2"
)

// ErrorCode type
type ErrorCode string

// AppError represents an application error with a code, message, cause, and extra data
type AppError struct {
	ErrorCode ErrorCode        // unique error code
	Message   string           // optional human-readable message
	Cause     error            // wrapped error
	ExtraData common.ExtraData // extra data associated with the error
}

func (ae *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", ae.GetCode(), ae.GetMessage())
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
		ae.ExtraData = make(common.ExtraData)
	}
	ae.ExtraData[key] = val
}

func (ae *AppError) GetExtraData() common.ExtraData {
	return ae.ExtraData
}

func (ae *AppError) Format(s fmt.State, verb rune) {
	errors.FormatError(ae, s, verb)
}

// Raise is a general util to construct domain errors
func Raise(ctx context.Context, code ErrorCode, msg string, errToWrap error, extraData common.ExtraData) error {
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
		ae.Cause = goerrors.Wrap(errToWrap, 1)
	} else {
		ae.Cause = goerrors.Wrap(fmt.Errorf("%s", ae.Error()), 1)
	}
	if ctx != nil {
		EnrichWithCtx(ctx, ae)
	}
	return ae
}

// RaiseInternal raises an internal error
func RaiseInternal(ctx context.Context, msg string, errToWrap error, extraData common.ExtraData) error {
	return Raise(ctx, ErrInternalError, msg, errToWrap, extraData)
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
			if e.Cause != nil {
				event.Extra["cause"] = e.Cause.Error()
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
				"context": common.GetExtraDataFromCtx(ctx),
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
