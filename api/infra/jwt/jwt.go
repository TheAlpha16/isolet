package jwt

import (
	"context"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type jwtImpl struct {
	key       []byte
	algorithm jwa.SignatureAlgorithm
}

func (jwtImpl *jwtImpl) Sign(ctx context.Context, claims *Claims) (string, error) {
	extraData := common.ExtraData{"claims": claims}

	jwtToken, err := claims.ToJWTToken(ctx)
	if err != nil {
		return "", errorDom.RaiseInternal(ctx, "", err, extraData)
	}
	extraData["jwt"] = jwtToken

	signedToken, err := jwt.Sign(jwtToken, jwt.WithKey(jwtImpl.algorithm, jwtImpl.key))
	if err != nil {
		return "", errorDom.Raise(ctx, errorDom.ErrTokenSigningFailed, "", err, extraData)
	}
	return string(signedToken), nil
}

func (jwtImpl *jwtImpl) Verify(ctx context.Context, token string) (*Claims, error) {
	extraData := common.ExtraData{"token": token}

	jwtToken, err := jwt.ParseString(token, jwt.WithKey(jwtImpl.algorithm, jwtImpl.key))
	if err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenExpiredInvalid, "", err, extraData)
	}

	claims, err := JWTTokenToClaims(ctx, jwtToken)
	if err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenExpiredInvalid, "", err, extraData)
	}

	return claims, nil
}

func NewJWT(key string) JWT {
	return &jwtImpl{
		key:       []byte(key),
		algorithm: jwa.HS256(),
	}
}
