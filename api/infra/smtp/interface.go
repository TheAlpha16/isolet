package smtp

import (
	"context"
	"time"

	emailDom "github.com/TheAlpha16/isolet/api/internal/domain/email"
)

type Creds struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Config struct {
	Creds
	Retries     int
	ConnTimeout time.Duration
	ChannelSize int
}

type SMTP interface {
	SendAsync(ctx context.Context, email *emailDom.Email) error
}
