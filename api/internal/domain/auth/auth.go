package auth

import (
	"time"

	"github.com/TheAlpha16/isolet/api/internal/domain"
)

type Session struct {
	ID        string
	UserID    int64
	ExpiresAt time.Time
	domain.BaseEntity
}
