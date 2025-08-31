package user

import (
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type UserUsecase struct {
	repo userDom.Repository
}

func (u *UserUsecase) Create(user *userDom.User) error {
	return u.repo.Create(user)
}

func (u *UserUsecase) Update(user *userDom.User) error {
	return u.repo.Update(user)
}

func New(repo userDom.Repository) userDom.Usecase {
	return &UserUsecase{
		repo: repo,
	}
}
