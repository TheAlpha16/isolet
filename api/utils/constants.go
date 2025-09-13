package utils

// base route groups
const (
	RouteAuth = "/auth"
)

// auth routes
const (
	RouteAuthRegister       = "/register"
	RouteAuthLogin          = "/login"
	RouteAuthVerify         = "/verify"
	RouteAuthForgotPassword = "/forgot-password"
	RouteAuthResetPassword  = "/reset-password"
)

// misc routes
const (
	RoutePing = "/ping"
)

// rest constants
const (
	AuthTokenCookieName = "token"
	TokenQueryKey       = "token"
)

// context keys
const (
	ContextKeyUserID = "user_id"
	ContextKeyRole   = "role"
)
