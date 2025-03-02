package router

import (
	"github.com/TheAlpha16/isolet/api/handler"
	"github.com/TheAlpha16/isolet/api/middleware"

	"github.com/gofiber/fiber/v2"
)

var (
	PING string = "/ping"

	AUTH_GROUP string = "/auth"
	LOGIN string = "/login"
	REGISTER string = "/register"
	VERIFY string = "/verify"
	FORGOT_PASSWORD string = "/forgot-password"
	RESET_PASSWORD string = "/reset-password"
	METADATA string = "/metadata"

	ONBOARD_GROUP = "/onboard"
	TEAM_CREATE string = "/team/create"
	TEAM_JOIN string = "/team/join"
	JOIN_TEAM_INVITE string = "/team/invite"

	API_GROUP string = "/api"
	IDENTIFY string = "/identify"
	LOGOUT string = "/logout"

	CHALLS string = "/challs"
	SUBMIT string = "/submit"
	UNLOCK_HINT string = "/hint/unlock"

	SCOREBOARD_GROUP string = "/scoreboard"
	SHOW_SCOREBOARD string = "/"
	TOP_SCORES string = "/top"

	LAUNCH string = "/launch"
	STOP string = "/stop"
	STATUS string = "/status"
	EXTEND string = "/extend"

	PROFILE_GROUP string = "/profile"
	GET_SELF_TEAM string = "/team/self"
	GET_INVITE_TOKEN string = "/team/invite"
)

func SetupRoutes(app *fiber.App) {
	app.Get(PING, handler.Ping)
	app.Use(middleware.ResolveIP())

	SetupAuthRoutes(app)
	SetupOnboardRoutes(app)
	SetupAPIRoutes(app)
}
