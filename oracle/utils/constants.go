package utils

// base route groups
const (
	RouteAuth      = "/auth"
	RouteTeam      = "/team"
	RouteEvent     = "/event"
	RouteProfile   = "/profile"
	RouteChallenge = "/challenge"
	RouteScore     = "/score"
	RouteInstance  = "/instance"
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

// profile routes
const (
	RouteProfileMe   = "/me"
	RouteProfileTeam = "/team"
)

// event routes
const (
	RouteEventInfo = "/info"
)

// challenge routes
const (
	RouteChallengeList       = "/"
	RouteChallengeSubmitFlag = "/submit"
	RouteHintUnlock          = "/hint/unlock"
)

// score routes
const (
	RouteScoreboard = "/"
	RouteScoreGraph = "/graph"
)

// instance routes
const (
	RouteInstanceList   = "/"
	RouteInstanceStart  = "/start"
	RouteInstanceStop   = "/stop"
	RouteInstanceExtend = "/extend"
)

// misc routes
const (
	RoutePing = "/ping"
)

// rest constants
const (
	AuthTokenCookieName = "token"
	TokenQueryKey       = "token"
	PageQueryKey        = "page"
	PageSizeQueryKey    = "page_size"
	IPHeaderKey         = "X-Real-IP"
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
	ContextKeyIP        = "ip"
)

// regular expressions
const (
	RegexNonAlphanumeric = `[^a-z0-9]+`
	RegexGenericFlag     = `^(.*\{)([^}]*)\}$`
)

// numeric constants
const (
	InstanceFlagSuffixLength = 8
)
