package user

import (
	domain_user "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type UserRepo struct{}

func (userRepo *UserRepo) Create(user *domain_user.User) error {
	// Implementation for creating a user
	return nil
}

func (userRepo *UserRepo) Update(user *domain_user.User) error {
	// Implementation for updating a user
	return nil
}

func New() domain_user.Repository {
	return &UserRepo{}
}
