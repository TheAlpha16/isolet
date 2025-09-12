package email

type Type string

const (
	TypeVerification  Type = "verification"
	TypePasswordReset Type = "password_reset"
)

type Entity struct {
	Name    string
	Address string
}

type Email struct {
	Sender     Entity
	Recipients []Entity
	Subject    string
	Body       string
}
