package errors

// SVC prefix for error codes
const (
	errExtraDataKey ErrCtxKey = "ctxErrExData"
	SVC             string    = "API"
)

// Error codes
const ()

// Map of error codes to user-facing messages
var msgMap = map[ErrorCode]string{}

// Map of server-side error codes that need to be filtered
var ServerSideErrors = map[ErrorCode]struct{}{}
