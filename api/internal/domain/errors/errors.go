package errors

import (
	"context"
	"fmt"
)

// ErrorCode type
type ErrorCode string

// ExtraData type
type ExtraData map[string]string

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

func (ae *AppError) AddExtraData(key string, val string) {
	if ae.ExtraData == nil {
		ae.ExtraData = make(map[string]string)
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
	return ae
}
