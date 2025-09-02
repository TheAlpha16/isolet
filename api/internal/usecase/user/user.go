package user

import (
	"context"

	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type userImpl struct {
	repo userDom.Repository
}

func (u *userImpl) Create(ctx context.Context, user *userDom.User) error {
	return u.repo.Create(ctx, user)
}

func (u *userImpl) Update(ctx context.Context, user *userDom.User) error {
	return u.repo.Update(ctx, user)
}

func (u *userImpl) ExistsByEmailOrUsername(ctx context.Context, email, username string) (*userDom.IdentifierExistence, error) {
	/*
		was thinking about caching both username and email while email verification
		but this poses issue - two generals problem
		one of the SET call fails, need rollback
		exploring lua script - allows for atomic operations
	*/
	return u.repo.CheckIdentifiers(ctx, email, username)
}

func New(repo userDom.Repository) userDom.Usecase {
	return &userImpl{
		repo: repo,
	}
}
