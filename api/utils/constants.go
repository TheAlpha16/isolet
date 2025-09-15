package utils

// base route groups
const (
	RouteAuth  = "/auth"
	RouteTeam  = "/team"
	RouteEvent = "/event"
)

// auth routes
const (
	RouteAuthRegister       = "/register"
	RouteAuthLogin          = "/login"
	RouteAuthVerify         = "/verify"
	RouteAuthForgotPassword = "/forgot-password"
	RouteAuthResetPassword  = "/reset-password"
	RouteAuthLogout         = "/logout"
)

// team routes
const (
	RouteTeamCreate = "/create"
	RouteTeamJoin   = "/join"
	RouteTeamInvite = "/invite"
)

// event routes
const (
	RouteEventInfo = "/info"
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
	ContextKeyUserID    = "user_id"
	ContextKeyTeamID    = "team_id"
	ContextKeyRole      = "role"
	ContextKeySessionID = "session_id"
)
