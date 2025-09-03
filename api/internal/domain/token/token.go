package token

import (
	"fmt"
	"time"

	"github.com/TheAlpha16/isolet/api/internal/domain"
	"github.com/TheAlpha16/isolet/api/utils/validator"

	govalidator "github.com/go-playground/validator/v10"
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

var validPurposes = map[TokenPurpose]struct{}{
	TokenAuth:              {},
	TokenPasswordReset:     {},
	TokenEmailVerification: {},
	TokenTeamInvite:        {},
}

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

func purposeValidator(fl govalidator.FieldLevel) bool {
	_, ok := validPurposes[TokenPurpose(fl.Field().String())]
	return ok
}

func init() {
	validator.RegisterValidation("token_purpose", purposeValidator)
}
