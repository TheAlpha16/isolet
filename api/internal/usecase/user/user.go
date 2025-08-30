package user

import (
	domain_user "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type UserUsecase struct {
	repo domain_user.Repository
}

func (u *UserUsecase) Create(user *domain_user.User) error {
	return u.repo.Create(user)
}

func (u *UserUsecase) Update(user *domain_user.User) error {
	return u.repo.Update(user)
}

func New(repo domain_user.Repository) domain_user.Usecase {
	return &UserUsecase{
		repo: repo,
	}
}
