package token

import (
	"github.com/TheAlpha16/isolet/api/infra/cache"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
)

type tokenImpl struct {
	cache cache.Cache
}

func New(cache cache.Cache) tokenDom.Usecase {
	return &tokenImpl{
		cache: cache,
	}
}
