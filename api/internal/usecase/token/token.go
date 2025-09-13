package token

import (
	"context"
	"fmt"

	"github.com/TheAlpha16/isolet/api/infra/cache"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"

	"github.com/vmihailenco/msgpack/v5"
)

type tokenImpl struct {
	cache cache.Cache
}

func (t *tokenImpl) Create(ctx context.Context, token *tokenDom.Token, extras map[string]string) error {
	key := token.Key()
	valBytes, err := msgpack.Marshal(token)
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrMarshalError, "failed to marshal token", err, nil)
	}

	if extras == nil {
		extras = map[string]string{}
	}
	extras[key] = string(valBytes)

	return t.cache.SetManyWithExpiry(ctx, extras, token.ExpiresAt)
}

func (t *tokenImpl) Fetch(ctx context.Context, id *tokenDom.TokenIdentifier) (*tokenDom.Token, error) {
	val, err := t.cache.Get(ctx, id.Key())
	if err != nil {
		if errorDom.IsSameError(err, errorDom.ErrCacheMiss) {
			return nil, errorDom.Raise(ctx, errorDom.ErrTokenExpiredInvalid, "", err, nil)
		}
		return nil, err
	}

	var token tokenDom.Token
	if err := msgpack.Unmarshal([]byte(val), &token); err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrUnmarshalError, "failed to unmarshal token", err, nil)
	}
	return &token, nil
}

func (t *tokenImpl) Delete(ctx context.Context, id *tokenDom.TokenIdentifier) error {
	return t.cache.Delete(ctx, id.Key())
}

func (t *tokenImpl) CountUserTokens(ctx context.Context, purpose tokenDom.TokenPurpose, userID int64) (int, error) {
	tokenIdentifier := &tokenDom.TokenIdentifier{
		ID:       "*",
		EntityID: fmt.Sprintf("%d", userID),
		Purpose:  purpose,
	}
	keys, err := t.cache.GetKeys(ctx, tokenIdentifier.Key())
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

func New(cache cache.Cache) tokenDom.Usecase {
	return &tokenImpl{
		cache: cache,
	}
}
