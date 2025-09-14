package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/TheAlpha16/isolet/api/internal/domain/errors"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/utils/validator"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

const (
	userIDClaim  string = "user_id"
	teamIDClaim  string = "team_id"
	roleClaim    string = "role"
	purposeClaim string = "purpose"
)

type Claims struct {
	JWTID     string                `validate:"required"`
	Subject   string                `validate:"required"`
	Purpose   tokenDom.TokenPurpose `validate:"required,token_purpose"`
	CreatedAt time.Time             `validate:"required"`
	ExpiresAt time.Time             `validate:"required"`
	UserID    *int64                `validate:"omitempty"`
	TeamID    *int64                `validate:"omitempty"`
	Role      *userDom.Role         `validate:"omitempty,role"`
}

func (claims *Claims) ToJWTToken(ctx context.Context) (jwt.Token, error) {
	err := claims.Validate(ctx)
	if err != nil {
		return nil, err
	}

	builder := jwt.NewBuilder().
		JwtID(claims.JWTID).
		Subject(claims.Subject).
		IssuedAt(claims.CreatedAt).
		Expiration(claims.ExpiresAt).
		Claim(purposeClaim, claims.Purpose)

	if claims.UserID != nil {
		builder = builder.Claim(userIDClaim, *claims.UserID)
	}
	if claims.TeamID != nil {
		builder = builder.Claim(teamIDClaim, *claims.TeamID)
	}
	if claims.Role != nil {
		builder = builder.Claim(roleClaim, *claims.Role)
	}

	return builder.Build()
}

func (claims *Claims) Validate(ctx context.Context) error {
	err := validator.Validate(ctx, claims)
	if err != nil {
		return errors.Raise(ctx, errors.ErrTokenMalformed, "validation failed", err, nil)
	}

	var userIdRequired, teamIdRequired, roleRequired bool

	switch claims.Purpose {
	case tokenDom.TokenAuth:
		userIdRequired = true
		roleRequired = true
	case tokenDom.TokenPasswordReset:
		userIdRequired = true
	case tokenDom.TokenTeamInvite:
		teamIdRequired = true
	}

	if userIdRequired && claims.UserID == nil {
		return errors.Raise(ctx, errors.ErrTokenMalformed, "user_id is missing", nil, nil)
	}
	if teamIdRequired && claims.TeamID == nil {
		return errors.Raise(ctx, errors.ErrTokenMalformed, "team_id is missing", nil, nil)
	}
	if roleRequired && claims.Role == nil {
		return errors.Raise(ctx, errors.ErrTokenMalformed, "role is missing", nil, nil)
	}

	return nil
}

func JWTTokenToClaims(ctx context.Context, token jwt.Token) (*Claims, error) {
	var claims Claims
	var roleString, purpose string
	var userIDFloat, teamIDFloat float64
	var ok bool

	if claims.JWTID, ok = token.JwtID(); !ok {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "jti is missing", nil, nil)
	}
	if claims.Subject, ok = token.Subject(); !ok {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "sub is missing", nil, nil)
	}
	if claims.ExpiresAt, ok = token.Expiration(); !ok {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "exp is missing", nil, nil)
	}
	if claims.CreatedAt, ok = token.IssuedAt(); !ok {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "iat is missing", nil, nil)
	}

	if err := token.Get(purposeClaim, &purpose); err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "purpose is missing", nil, nil)
	}
	claims.Purpose = tokenDom.TokenPurpose(purpose)

	// These are optional - only parse if present
	token.Get(userIDClaim, &userIDFloat)
	token.Get(teamIDClaim, &teamIDFloat)
	token.Get(roleClaim, &roleString)

	if userIDFloat != 0 {
		userID := int64(userIDFloat)
		claims.UserID = &userID
	}
	if teamIDFloat != 0 {
		teamID := int64(teamIDFloat)
		claims.TeamID = &teamID
	}
	if roleString != "" {
		userRole := userDom.Role(roleString)
		claims.Role = &userRole
	}

	return &claims, claims.Validate(ctx)
}

func NewAuthClaims(jwtID string, userID int64, teamID *int64, role userDom.Role, expiresAt time.Time) *Claims {
	return &Claims{
		JWTID:     jwtID,
		Subject:   fmt.Sprintf("%d", userID),
		UserID:    &userID,
		TeamID:    teamID,
		Role:      &role,
		Purpose:   tokenDom.TokenAuth,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
}

func NewEmailVerificationClaims(jwtID, verificationID string, expiresAt time.Time) *Claims {
	return &Claims{
		JWTID:     jwtID,
		Subject:   verificationID,
		UserID:    nil,
		TeamID:    nil,
		Role:      nil,
		Purpose:   tokenDom.TokenEmailVerification,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
}

func NewPasswordResetClaims(jwtID string, userID int64, expiresAt time.Time) *Claims {
	return &Claims{
		JWTID:     jwtID,
		Subject:   fmt.Sprintf("%d", userID),
		UserID:    &userID,
		TeamID:    nil,
		Role:      nil,
		Purpose:   tokenDom.TokenPasswordReset,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
}

func NewTeamInviteClaims(jwtID string, teamID int64, expiresAt time.Time) *Claims {
	return &Claims{
		JWTID:     jwtID,
		Subject:   fmt.Sprintf("%d", teamID),
		UserID:    nil,
		TeamID:    &teamID,
		Role:      nil,
		Purpose:   tokenDom.TokenTeamInvite,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
}
