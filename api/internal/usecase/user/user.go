package user

import (
	"context"

	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type UserUsecase struct {
	repo userDom.Repository
}

func (u *UserUsecase) Create(ctx context.Context, user *userDom.User) error {
	return u.repo.Create(ctx, user)
}

func (u *UserUsecase) Update(ctx context.Context, user *userDom.User) error {
	return u.repo.Update(ctx, user)
}

func New(repo userDom.Repository) userDom.Usecase {
	return &UserUsecase{
		repo: repo,
	}
}
