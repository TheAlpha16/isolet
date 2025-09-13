package configvars

type ConfigKey[T any] struct {
	Name    string
	Default T
}

var (
	SMTPHost     = ConfigKey[string]{Name: "smtp_host", Default: "smtp.doesntexist.com"}
	SMTPPort     = ConfigKey[int]{Name: "smtp_port", Default: 587}
	SMTPUser     = ConfigKey[string]{Name: "smtp_user", Default: "user@doesntexist.com"}
	SMTPPassword = ConfigKey[string]{Name: "smtp_password", Default: "password"}
)
