package errors

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	goerrors "github.com/go-errors/errors"
)

// ErrorCode type
type ErrorCode string

// AppError represents an application error with a code, message, cause, and extra data
type AppError struct {
	ErrorCode ErrorCode // unique error code
	Message   string    // optional human-readable message
	Cause     error     // wrapped error
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

// Raise is a general util to construct domain errors
func Raise(code ErrorCode, msg string, errToWrap error) error {
	ae := &AppError{
		ErrorCode: code,
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
			event.Message = fmt.Sprintf("%s : %s", e.GetCode(), e.GetMessage())
			event.Exception = []sentry.Exception{{
				Value:      fmt.Sprintf("%s", e.GetCode()),
				Type:       fmt.Sprintf("%T", e),
				Stacktrace: sentry.ExtractStacktrace(e.Cause),
			}}
			if e.Cause != nil {
				event.Extra["cause"] = e.Cause.Error()
			}
		})
		hub.CaptureEvent(event)
	default:
		stackErr := goerrors.Wrap(e, 1)
		hub.WithScope(func(scope *sentry.Scope) {
			event.Message = fmt.Sprintf("%s", e.Error())
			event.Exception = []sentry.Exception{{
				Value:      e.Error(),
				Type:       fmt.Sprintf("%T", e),
				Stacktrace: sentry.ExtractStacktrace(stackErr),
			}}
		})
		addTraceContextToSentryEvent(ctx, event)
		hub.CaptureEvent(event)
	}
}
