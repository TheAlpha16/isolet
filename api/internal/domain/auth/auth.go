package auth

import (
	"time"
)

type Session struct {
	UserID    int64     `json:"user_id"`
	Token     string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
}
