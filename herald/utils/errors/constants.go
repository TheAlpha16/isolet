package errors

const (
	SVC string = "HRLD"
)

// Error codes
const (
	// System errors
	ErrInternalError  ErrorCode = "SYS-00"
)

// Map of error codes to user-facing messages
var msgMap = map[ErrorCode]string{
	// System errors
	ErrInternalError: "internal server error",
}
