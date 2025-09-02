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

func New(repo userDom.Repository) userDom.Usecase {
	return &userImpl{
		repo: repo,
	}
}
