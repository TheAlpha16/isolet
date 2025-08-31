package errors

// SVC prefix for error codes
const (
	errExtraDataKey ErrCtxKey = "ctxErrExData"
	SVC             string    = "API"
)

// Error codes
const (
	ErrValidationFailed ErrorCode = "VAL-00"
)

// Map of error codes to user-facing messages
var msgMap = map[ErrorCode]string{
	ErrValidationFailed: "one or more validation errors occurred",
}

// Map of server-side error codes that need to be filtered
var ServerSideErrors = map[ErrorCode]struct{}{}
