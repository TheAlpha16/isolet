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
	// System errors
	ErrInternalServerError ErrorCode = "SYS-00"

	// Validation errors
	ErrValidationFailed ErrorCode = "VAL-00"

	// Rest errors
	ErrInvalidPayload ErrorCode = "REST-00"

	// User errors
	ErrInvalidRole ErrorCode = "USER-00"

	// DB errors
	ErrDBCreateError ErrorCode = "DB-00"
)

// Map of error codes to user-facing messages
var msgMap = map[ErrorCode]string{
	// System errors
	ErrInternalServerError: "internal server error",

	// Validation errors
	ErrValidationFailed: "one or more validation errors occurred",

	// Rest errors
	ErrInvalidPayload: "invalid request payload",

	// User errors
	ErrInvalidRole: "invalid user role",
}

// Map of server-side error codes that need to be filtered
var ServerSideErrors = map[ErrorCode]struct{}{
	ErrInternalServerError: {},
}

// Map of auth error codes
var AuthErrors = map[ErrorCode]struct{}{}

// Map of forbidden error codes
var ForbiddenErrors = map[ErrorCode]struct{}{}
