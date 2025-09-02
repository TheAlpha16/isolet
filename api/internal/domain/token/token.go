package token

import (
	"time"

	"github.com/TheAlpha16/isolet/api/internal/domain"
)

type TokenPurpose string

const (
	TokenAuth              TokenPurpose = "auth"
	TokenPasswordReset     TokenPurpose = "password_reset"
	TokenEmailVerification TokenPurpose = "email_verification"
	TokenTeamInvite        TokenPurpose = "team_invite"
)

type Token struct {
	ID        string
	EntityID  string
	Metadata  map[string]string
	ExpiresAt time.Time
	domain.BaseEntity
}
