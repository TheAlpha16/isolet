package configvars

type ConfigKey[T any] struct {
	Name    string
	Default T
}

var (
	// general
	EventName = ConfigKey[string]{Name: "event_name", Default: "isolet"}
	PublicURL = ConfigKey[string]{Name: "public_url", Default: "https://isolet.dev"}

	//smtp
	SMTPHost     = ConfigKey[string]{Name: "smtp_host", Default: "smtp.isolet.dev"}
	SMTPPort     = ConfigKey[int]{Name: "smtp_port", Default: 587}
	SMTPUser     = ConfigKey[string]{Name: "smtp_user", Default: "user@isolet.dev"}
	SMTPPassword = ConfigKey[string]{Name: "smtp_password", Default: "password"}

	// email
	EmailVerificationEnabled = ConfigKey[bool]{Name: "email_verification_enabled", Default: true}
	EmailSender              = ConfigKey[string]{Name: "email_from", Default: "noreply@isolet.dev"}

	// auth
	AuthMaxSessions          = ConfigKey[int]{Name: "auth_max_sessions", Default: 5}
	PasswordResetEnabled     = ConfigKey[bool]{Name: "password_reset_enabled", Default: true}
	PasswordResetMaxSessions = ConfigKey[int]{Name: "password_reset_max_sessions", Default: 1}

	// profile
	TeamMaxSize = ConfigKey[int]{Name: "team_max_size", Default: 4}
)
