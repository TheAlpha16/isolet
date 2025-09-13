package user

import (
	"context"

	"github.com/TheAlpha16/isolet/api/infra/cache"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type userImpl struct {
	repo  userDom.Repository
	cache cache.Cache
}

func (u *userImpl) Create(ctx context.Context, user *userDom.User) (*userDom.User, error) {
	return u.repo.Create(ctx, user)
}

func (u *userImpl) Update(ctx context.Context, user *userDom.User) error {
	return u.repo.Update(ctx, user)
}

func (u *userImpl) GetByUsernameOrEmail(ctx context.Context, identifier string) (*userDom.User, error) {
	return u.repo.GetByUsernameOrEmail(ctx, identifier)
}

func (u *userImpl) ExistsByEmailOrUsername(ctx context.Context, email, username string) (*userDom.IdentifierExistence, error) {
	exists := userDom.IdentifierExistence{
		EmailExists:    false,
		UsernameExists: false,
	}

	_, err := u.cache.Get(ctx, userDom.EmailCacheKey(email))
	if err != nil {
		// return immediately in case of errors other than cache miss
		if !errorDom.IsSameError(err, errorDom.ErrCacheMiss) {
			return nil, err
		}
	} else {
		exists.EmailExists = true
		return &exists, nil
	}

	_, err = u.cache.Get(ctx, userDom.UsernameCacheKey(username))
	if err != nil {
		// return immediately in case of errors other than cache miss
		if !errorDom.IsSameError(err, errorDom.ErrCacheMiss) {
			return nil, err
		}
	} else {
		exists.UsernameExists = true
		return &exists, nil
	}

	return u.repo.CheckIdentifiers(ctx, email, username)
}

func New(cache cache.Cache, repo userDom.Repository) userDom.Usecase {
	return &userImpl{
		repo:  repo,
		cache: cache,
	}
}
