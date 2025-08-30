package errors

import "errors"

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
