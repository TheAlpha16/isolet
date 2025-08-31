package errors

const (
	// Context key for storing/retrieving error extra data
	errExtraDataKey ErrCtxKey = "ctxErrExData"
	srcKey          ErrCtxKey = "ctxErrSrc"

	// SVC prefix for error codes
	SVC string = "API"
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
