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

func (u *userImpl) GetByEmailOrUsername(ctx context.Context, email, username string) (*userDom.User, error) {
	return u.repo.GetByEmailOrUsername(ctx, email, username)
}

func (u *userImpl) ExistsByEmailOrUsername(ctx context.Context, email, username string) error {
	emailTaken := errorDom.Raise(ctx, errorDom.ErrUserEmailTaken, "", nil, nil)
	usernameTaken := errorDom.Raise(ctx, errorDom.ErrUserUsernameTaken, "", nil, nil)

	_, err := u.cache.Get(ctx, userDom.EmailCacheKey(email))
	if err != nil {
		// return immediately in case of errors other than cache miss
		if !errorDom.IsSameError(err, errorDom.ErrCacheMiss) {
			return err
		}
	} else {
		return emailTaken
	}

	_, err = u.cache.Get(ctx, userDom.UsernameCacheKey(username))
	if err != nil {
		// return immediately in case of errors other than cache miss
		if !errorDom.IsSameError(err, errorDom.ErrCacheMiss) {
			return err
		}
	} else {
		return usernameTaken
	}

	user, err := u.repo.GetByEmailOrUsername(ctx, email, username)
	if err != nil {
		if errorDom.IsSameError(err, errorDom.ErrUserNotFound) {
			return nil
		}
		return err
	}

	if user.Email == email {
		return emailTaken
	}

	return usernameTaken
}

func New(cache cache.Cache, repo userDom.Repository) userDom.Usecase {
	return &userImpl{
		repo:  repo,
		cache: cache,
	}
}
