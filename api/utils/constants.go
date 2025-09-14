package utils

// base route groups
const (
	RouteAuth = "/auth"
	RouteTeam = "/team"
)

// auth routes
const (
	RouteAuthRegister       = "/register"
	RouteAuthLogin          = "/login"
	RouteAuthVerify         = "/verify"
	RouteAuthForgotPassword = "/forgot-password"
	RouteAuthResetPassword  = "/reset-password"
)

// team routes
const (
	RouteTeamCreate = "/create"
	RouteTeamJoin   = "/join"
	RouteTeamInvite = "/invite"
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

// front-end routes
const (
	RouteFrontResetPassword = "/reset-password"
)

// context keys
const (
	ContextKeyUserID = "user_id"
	ContextKeyTeamID = "team_id"
	ContextKeyRole   = "role"
)
