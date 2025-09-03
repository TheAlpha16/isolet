package jwt

import (
	"context"
	"time"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/utils/validator"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

const (
	userIDClaim  string = "user_id"
	roleClaim    string = "role"
	purposeClaim string = "purpose"
)

type Claims struct {
	JWTID     string                `validate:"required"`
	UserID    int64                 `validate:"required"`
	Role      userDom.Role          `validate:"required,role"`
	Purpose   tokenDom.TokenPurpose `validate:"required,token_purpose"`
	CreatedAt time.Time             `validate:"required"`
	ExpiresAt time.Time             `validate:"required"`
}

func (claims *Claims) ToJWTToken(ctx context.Context) (jwt.Token, error) {
	err := claims.Validate(ctx)
	if err != nil {
		return nil, err
	}

	builder := jwt.NewBuilder().
		JwtID(claims.JWTID).
		IssuedAt(claims.CreatedAt).
		Expiration(claims.ExpiresAt).
		Claim(userIDClaim, claims.UserID).
		Claim(roleClaim, claims.Role).
		Claim(purposeClaim, claims.Purpose)

	return builder.Build()
}

func (claims *Claims) Validate(ctx context.Context) error {
	return validator.Validate(ctx, claims)
}

func JWTTokenToClaims(ctx context.Context, token jwt.Token) (*Claims, error) {
	var claims Claims
	var err error
	var ok bool

	if claims.JWTID, ok = token.JwtID(); !ok {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "jti is missing", nil, nil)
	}
	if claims.ExpiresAt, ok = token.Expiration(); !ok {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "exp is missing", nil, nil)
	}
	if claims.CreatedAt, ok = token.IssuedAt(); !ok {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "iat is missing", nil, nil)
	}
	if err = token.Get(userIDClaim, &claims.UserID); err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "user_id is missing", nil, nil)
	}
	if err = token.Get(roleClaim, &claims.Role); err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "role is missing", nil, nil)
	}
	if err = token.Get(purposeClaim, &claims.Purpose); err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "purpose is missing", nil, nil)
	}

	if err = claims.Validate(ctx); err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "claims are invalid", err, nil)
	}
	return &claims, nil
}
