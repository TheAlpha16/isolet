package errors

const (
	SVC string = "HRLD"
)

// Error codes
const (
	// System errors
	ErrInternalError ErrorCode = "SYS-00"

	// Fact errors
	ErrFactInvalidKey ErrorCode = "FACT-01"
)

// Map of error codes to user-facing messages
var msgMap = map[ErrorCode]string{
	// System errors
	ErrInternalError: "internal server error",

	// Fact errors
	ErrFactInvalidKey: "fact has invalid key",
}
