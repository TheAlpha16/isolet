package errors

import "github.com/TheAlpha16/isolet/api/internal/domain/common"

const (
	// Context key for storing/retrieving error extra data
	srcKey common.ContextKey = "ctxErrSrc"

	// SVC prefix for error codes
	SVC string = "API"
)

// Error codes
const (
	// System errors
	ErrInternalError  ErrorCode = "SYS-00"
	ErrMarshalError   ErrorCode = "SYS-01"
	ErrUnmarshalError ErrorCode = "SYS-02"

	// Validation errors
	ErrValidationFailed ErrorCode = "VAL-00"

	// Rest errors
	ErrRestInvalidPayload ErrorCode = "REST-00"

	// User errors
	ErrUserInvalidRole   ErrorCode = "USER-00"
	ErrUserEmailTaken    ErrorCode = "USER-01"
	ErrUserUsernameTaken ErrorCode = "USER-02"
	ErrUserNotFound      ErrorCode = "USER-03"

	// DB errors
	ErrDBCreateError ErrorCode = "DB-00"
	ErrDBReadError   ErrorCode = "DB-01"

	// Cache errors
	ErrCacheCallFail       ErrorCode = "CACHE-00"
	ErrCacheMiss           ErrorCode = "CACHE-01"
	ErrCacheScriptLoadFail ErrorCode = "CACHE-02"

	// Token errors
	ErrTokenMalformed      ErrorCode = "TOKEN-00"
	ErrTokenExpiredInvalid ErrorCode = "TOKEN-01"
	ErrTokenSigningFailed  ErrorCode = "TOKEN-02"

	// ConfigVar errors
	ErrConfigVarInvalid       ErrorCode = "CONFIGVAR-00"
	ErrConfigVarRefreshFailed ErrorCode = "CONFIGVAR-01"

	// SMTP errors
	ErrSMTPDialError          ErrorCode = "SMTP-00"
	ErrSMTPQueueFull          ErrorCode = "SMTP-01"
	ErrSMTPMaxRetriesExceeded ErrorCode = "SMTP-02"
	ErrSMTPSendFailed         ErrorCode = "SMTP-03"
	ErrSMTPContextCanceled    ErrorCode = "SMTP-04"

	// Email errors
	ErrEmailInvalidType     ErrorCode = "EMAIL-00"
	ErrEmailTemplateFetch   ErrorCode = "EMAIL-01"
	ErrEmailTemplateExecute ErrorCode = "EMAIL-02"
	ErrEmailLinkBuild       ErrorCode = "EMAIL-03"

	// Auth errors
	ErrAuthInvalidCredentials ErrorCode = "AUTH-00"
	ErrAuthMaxSessionsReached ErrorCode = "AUTH-01"
)

// Map of error codes to user-facing messages
var msgMap = map[ErrorCode]string{
	// System errors
	ErrInternalError: "internal server error",

	// Validation errors
	ErrValidationFailed: "one or more validation errors occurred",

	// Rest errors
	ErrRestInvalidPayload: "invalid request payload",

	// User errors
	ErrUserInvalidRole:   "invalid user role",
	ErrUserEmailTaken:    "email is already taken",
	ErrUserUsernameTaken: "username is already taken",
	ErrUserNotFound:      "user not found",

	// Cache errors
	ErrCacheMiss:           "cache miss",
	ErrCacheScriptLoadFail: "failed to load cache script",

	// Token errors
	ErrTokenExpiredInvalid: "token is invalid or expired",
	ErrTokenSigningFailed:  "token signing failed",

	// ConfigVar errors
	ErrConfigVarInvalid:       "config variable is invalid",
	ErrConfigVarRefreshFailed: "config variable refresh failed",

	// SMTP errors
	ErrSMTPDialError:          "failed to dial SMTP server",
	ErrSMTPQueueFull:          "SMTP queue is full",
	ErrSMTPMaxRetriesExceeded: "email max retries exceeded",
	ErrSMTPSendFailed:         "failed to send email",
	ErrSMTPContextCanceled:    "SMTP context canceled",

	// Email errors
	ErrEmailInvalidType:     "invalid email type",
	ErrEmailTemplateFetch:   "failed to fetch email template",
	ErrEmailTemplateExecute: "failed to execute email template",
	ErrEmailLinkBuild:       "failed to build email link",

	// Auth errors
	ErrAuthInvalidCredentials: "invalid credentials",
	ErrAuthMaxSessionsReached: "maximum sessions reached",
}

// Map of server-side error codes that need to be filtered
var ServerSideErrors = map[ErrorCode]struct{}{
	ErrInternalError:          {},
	ErrCacheCallFail:          {},
	ErrMarshalError:           {},
	ErrUnmarshalError:         {},
	ErrDBCreateError:          {},
	ErrDBReadError:            {},
	ErrCacheScriptLoadFail:    {},
	ErrTokenSigningFailed:     {},
	ErrConfigVarInvalid:       {},
	ErrConfigVarRefreshFailed: {},
	ErrSMTPQueueFull:          {},
	ErrSMTPContextCanceled:    {},
	ErrEmailInvalidType:       {},
	ErrEmailTemplateFetch:     {},
	ErrEmailTemplateExecute:   {},
	ErrEmailLinkBuild:         {},
}

// Map of auth error codes
var AuthErrors = map[ErrorCode]struct{}{
	ErrTokenExpiredInvalid:    {},
	ErrAuthInvalidCredentials: {},
}

// Map of forbidden error codes
var ForbiddenErrors = map[ErrorCode]struct{}{
	ErrAuthMaxSessionsReached: {},
}
