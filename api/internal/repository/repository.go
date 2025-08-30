package repository

import (
	domain_user "github.com/TheAlpha16/isolet/api/internal/domain/user"

	repo_user "github.com/TheAlpha16/isolet/api/internal/repository/user"
)

type Repositories struct {
	User domain_user.Repository
}

func New() *Repositories {
	userRepo := repo_user.New()

	return &Repositories{
		User: userRepo,
	}
}
