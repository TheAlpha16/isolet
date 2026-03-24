package errors

import (
	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
)

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
	ErrRestMissingToken   ErrorCode = "REST-01"
	ErrRestTimedOut       ErrorCode = "REST-02"

	// User errors
	ErrUserInvalidRole     ErrorCode = "USER-00"
	ErrUserEmailTaken      ErrorCode = "USER-01"
	ErrUserUsernameTaken   ErrorCode = "USER-02"
	ErrUserNotFound        ErrorCode = "USER-03"
	ErrUserInsufficentRole ErrorCode = "USER-04"

	// Team errors
	ErrTeamRequired      ErrorCode = "TEAM-00"
	ErrTeamAlreadyJoined ErrorCode = "TEAM-01"
	ErrTeamNameTaken     ErrorCode = "TEAM-02"
	ErrTeamNotFound      ErrorCode = "TEAM-03"
	ErrTeamFull          ErrorCode = "TEAM-04"

	// DB errors
	ErrDBCreateError ErrorCode = "DB-00"
	ErrDBReadError   ErrorCode = "DB-01"
	ErrDBUpdateError ErrorCode = "DB-02"
	ErrDBExecError   ErrorCode = "DB-03"
	ErrDBDeleteError ErrorCode = "DB-04"

	// Cache errors
	ErrCacheCallFail          ErrorCode = "CACHE-00"
	ErrCacheMiss              ErrorCode = "CACHE-01"
	ErrCacheScriptLoadFail    ErrorCode = "CACHE-02"
	ErrCacheZSetMissingMember ErrorCode = "CACHE-03"

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
	ErrAuthInvalidCredentials    ErrorCode = "AUTH-00"
	ErrAuthMaxSessionsReached    ErrorCode = "AUTH-01"
	ErrAuthPasswordResetDisabled ErrorCode = "AUTH-02"
	ErrAuthMissingToken          ErrorCode = "AUTH-03"
	ErrAuthUserBanned            ErrorCode = "AUTH-04"

	// Challenge errors
	ErrChallengeNotFound           ErrorCode = "CHALLENGE-00"
	ErrChallengeAlreadySolved      ErrorCode = "CHALLENGE-01"
	ErrChallengeMaxAttemptsReached ErrorCode = "CHALLENGE-02"
	ErrChallengeNameInvalid        ErrorCode = "CHALLENGE-03"

	// Hint errors
	ErrHintNotFound        ErrorCode = "HINT-00"
	ErrHintCostExceeded    ErrorCode = "HINT-01"
	ErrHintAlreadyUnlocked ErrorCode = "HINT-02"

	// Scoreboard errors
	ErrScoreUpdateFailed ErrorCode = "SCORE-00"

	// Instance errors
	ErrInstanceNotOnDemand             ErrorCode = "INSTANCE-00"
	ErrInstanceNotFound                ErrorCode = "INSTANCE-01"
	ErrInstanceWIP                     ErrorCode = "INSTANCE-02"
	ErrInstanceAlreadyRunning          ErrorCode = "INSTANCE-03"
	ErrInstanceCreationFailed          ErrorCode = "INSTANCE-04"
	ErrInstanceInvalid                 ErrorCode = "INSTANCE-05"
	ErrInstanceInvalidResourceQuantity ErrorCode = "INSTANCE-06"
	ErrInstanceExtensionNotAllowed     ErrorCode = "INSTANCE-07"
	ErrInstanceUpdateFailed            ErrorCode = "INSTANCE-08"
	ErrInstanceDeletionFailed          ErrorCode = "INSTANCE-09"
	ErrInstanceInvalidState            ErrorCode = "INSTANCE-10"
	ErrInstanceNotRunning              ErrorCode = "INSTANCE-11"

	// k8s errors
	ErrK8sConnectionFailed ErrorCode = "K8S-00"
	ErrK8sInstanceNotFound ErrorCode = "K8S-01"

	// Manifest errors
	ErrManifestNotFound ErrorCode = "MANIFEST-00"

	// Kafka errors
	ErrKafkaConnectionFailed   ErrorCode = "KAFKA-00"
	ErrKafkaSubscriptionFailed ErrorCode = "KAFKA-01"
	ErrKafkaConsumerError      ErrorCode = "KAFKA-02"

	// Fact errors
	ErrFactDeserialization ErrorCode = "FACT-00"

	// Event errors
	ErrEventNotStarted ErrorCode = "EVENT-00"
	ErrEventEnded      ErrorCode = "EVENT-01"
)

// Map of error codes to user-facing messages
var msgMap = map[ErrorCode]string{
	// System errors
	ErrInternalError: "internal server error",

	// Validation errors
	ErrValidationFailed: "one or more validation errors occurred",

	// Rest errors
	ErrRestInvalidPayload: "invalid request payload",
	ErrRestMissingToken:   "token required for verification",
	ErrRestTimedOut:       "request timed out, please try again later",

	// User errors
	ErrUserInvalidRole:     "invalid user role",
	ErrUserEmailTaken:      "email is already taken",
	ErrUserUsernameTaken:   "username is already taken",
	ErrUserNotFound:        "user not found",
	ErrUserInsufficentRole: "insufficient role",

	// Team errors
	ErrTeamRequired:      "need to be a member of a team",
	ErrTeamAlreadyJoined: "user is already in a team",
	ErrTeamNameTaken:     "team name is already taken",
	ErrTeamNotFound:      "team not found",
	ErrTeamFull:          "team is full",

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
	ErrAuthInvalidCredentials:    "invalid credentials",
	ErrAuthMaxSessionsReached:    "maximum sessions reached",
	ErrAuthPasswordResetDisabled: "password reset is disabled",
	ErrAuthMissingToken:          "missing auth token",
	ErrAuthUserBanned:            "user is banned",

	// Challenge errors
	ErrChallengeNotFound:           "challenge not found",
	ErrChallengeAlreadySolved:      "challenge already solved",
	ErrChallengeMaxAttemptsReached: "maximum attempts reached",
	ErrChallengeNameInvalid:        "challenge name is invalid",

	// Hint errors
	ErrHintNotFound:        "hint not found",
	ErrHintCostExceeded:    "insufficient points to unlock hint",
	ErrHintAlreadyUnlocked: "hint already unlocked",

	// Score errors
	ErrScoreUpdateFailed: "failed to update score",

	// Instance errors
	ErrInstanceNotOnDemand:             "instances are not available for this challenge",
	ErrInstanceNotFound:                "instance not found",
	ErrInstanceWIP:                     "instance work in progress, please wait a few seconds before retrying",
	ErrInstanceAlreadyRunning:          "instance is already running",
	ErrInstanceCreationFailed:          "failed to create instance",
	ErrInstanceInvalidResourceQuantity: "instance has invalid resource quantity",
	ErrInstanceExtensionNotAllowed:     "instance extension is not allowed",
	ErrInstanceInvalidState:            "instance is in an invalid state",
	ErrInstanceNotRunning:              "instance is not running",

	// k8s errors
	ErrK8sInstanceNotFound: "instance not found in kubernetes cluster",

	// Manifest errors
	ErrManifestNotFound: "manifest not found",

	// Kafka errors
	ErrKafkaConnectionFailed:   "failed to connect to kafka",
	ErrKafkaSubscriptionFailed: "failed to subscribe to kafka topics",
	ErrKafkaConsumerError:      "kafka consumer encountered an error",

	// Fact errors
	ErrFactDeserialization: "failed to deserialize fact",

	// Event errors
	ErrEventNotStarted: "event has not started yet",
	ErrEventEnded:      "event has already ended",
}

// Map of server-side error codes that need to be filtered
var ServerErrorCodes = map[ErrorCode]struct{}{
	ErrInternalError:           {},
	ErrCacheCallFail:           {},
	ErrMarshalError:            {},
	ErrUnmarshalError:          {},
	ErrDBCreateError:           {},
	ErrDBReadError:             {},
	ErrDBUpdateError:           {},
	ErrDBExecError:             {},
	ErrDBDeleteError:           {},
	ErrCacheScriptLoadFail:     {},
	ErrTokenSigningFailed:      {},
	ErrConfigVarInvalid:        {},
	ErrConfigVarRefreshFailed:  {},
	ErrSMTPQueueFull:           {},
	ErrSMTPContextCanceled:     {},
	ErrEmailInvalidType:        {},
	ErrEmailTemplateFetch:      {},
	ErrEmailTemplateExecute:    {},
	ErrEmailLinkBuild:          {},
	ErrScoreUpdateFailed:       {},
	ErrCacheZSetMissingMember:  {},
	ErrInstanceCreationFailed:  {},
	ErrK8sConnectionFailed:     {},
	ErrK8sInstanceNotFound:     {},
	ErrInstanceInvalid:         {},
	ErrChallengeNameInvalid:    {},
	ErrInstanceUpdateFailed:    {},
	ErrManifestNotFound:        {},
	ErrInstanceDeletionFailed:  {},
	ErrKafkaConnectionFailed:   {},
	ErrKafkaSubscriptionFailed: {},
	ErrKafkaConsumerError:      {},
	ErrFactDeserialization:     {},
}

// Map of auth error codes
var AuthErrorCodes = map[ErrorCode]struct{}{
	ErrTokenExpiredInvalid:    {},
	ErrAuthInvalidCredentials: {},
	ErrAuthMissingToken:       {},
}

// Map of forbidden error codes
var ForbiddenErrorCodes = map[ErrorCode]struct{}{
	ErrAuthMaxSessionsReached:    {},
	ErrAuthPasswordResetDisabled: {},
	ErrUserInsufficentRole:       {},
	ErrTeamRequired:              {},
	ErrAuthUserBanned:            {},
	ErrEventNotStarted:           {},
	ErrEventEnded:                {},
}

// postgres error codes
const (
	PgErrDuplicateKey = "23505"
)
