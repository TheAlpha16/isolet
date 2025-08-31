package user

import (
	"context"

	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func (userRepo *UserRepo) Create(ctx context.Context, user *userDom.User) error {
	// Implementation for creating a user
	return nil
}

func (userRepo *UserRepo) Update(ctx context.Context, user *userDom.User) error {
	// Implementation for updating a user
	return nil
}

func New(db *gorm.DB) userDom.Repository {
	return &UserRepo{
		db: db,
	}
}
