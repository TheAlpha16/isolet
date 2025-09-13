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
	EmailVerification = ConfigKey[bool]{Name: "email_verification", Default: true}
	EmailSender       = ConfigKey[string]{Name: "email_from", Default: "noreply@isolet.dev"}
)
