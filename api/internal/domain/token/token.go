package token

import (
	"fmt"
	"time"

	"github.com/TheAlpha16/isolet/api/internal/domain"
)

type TokenPurpose string

const (
	tokenPrefix string = "token"
)

const (
	TokenAuth              TokenPurpose = "auth"
	TokenPasswordReset     TokenPurpose = "password_reset"
	TokenEmailVerification TokenPurpose = "email_verification"
	TokenTeamInvite        TokenPurpose = "team_invite"
)

type TokenIdentifier struct {
	ID       string
	EntityID string
	Purpose  TokenPurpose
}

type Token struct {
	TokenIdentifier
	Metadata  map[string]string `msgpack:"metadata"`
	ExpiresAt time.Time         `msgpack:"expires_at"`
	domain.BaseEntity
}

func (tid *TokenIdentifier) Key() string {
	return fmt.Sprintf("%s:%s:%s:%s", tokenPrefix, tid.Purpose, tid.EntityID, tid.ID)
}
