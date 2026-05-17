package configvars

type ConfigKey[T any] struct {
	Name    string
	Default T
}

var (
	// event
	EventName      = ConfigKey[string]{Name: "event_name", Default: "isolet"}
	EventPublicURL = ConfigKey[string]{Name: "event_public_url", Default: "https://isolet.dev"}
	EventStart     = ConfigKey[int]{Name: "event_start", Default: 1709899200}
	EventEnd       = ConfigKey[int]{Name: "event_end", Default: 1710028800}
	EventPostMode  = ConfigKey[bool]{Name: "event_post_mode", Default: false}

	// smtp
	SMTPHost     = ConfigKey[string]{Name: "smtp_host", Default: "smtp.isolet.dev"}
	SMTPPort     = ConfigKey[int]{Name: "smtp_port", Default: 587}
	SMTPUser     = ConfigKey[string]{Name: "smtp_user", Default: "user@isolet.dev"}
	SMTPPassword = ConfigKey[string]{Name: "smtp_password", Default: "password"}

	// email
	EmailVerificationEnabled = ConfigKey[bool]{Name: "email_verification_enabled", Default: true}
	EmailSender              = ConfigKey[string]{Name: "email_from", Default: "noreply@isolet.dev"}

	// auth
	AuthMaxSessions              = ConfigKey[int]{Name: "auth_max_sessions", Default: 5}
	AuthPasswordResetEnabled     = ConfigKey[bool]{Name: "auth_password_reset_enabled", Default: true}
	AuthPasswordResetMaxSessions = ConfigKey[int]{Name: "auth_password_reset_max_sessions", Default: 1}

	// team
	TeamMaxSize = ConfigKey[int]{Name: "team_max_size", Default: 4}

	// instance
	InstanceDomain = ConfigKey[string]{Name: "instance_domain", Default: "isolet.dev"}
)
