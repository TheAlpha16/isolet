package email

type EmailIdentifier struct {
	Username string
	To       string
}

type EmailInput struct {
	EmailIdentifier
	Type  Type
	Token string
}

type TemplateInput struct {
	EventName string
	Username  string
	Link      string
}
