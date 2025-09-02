package auth

import (
	"time"
)

type Session struct {
	UserID    int64
	Token     string
	ExpiresAt time.Time
}
