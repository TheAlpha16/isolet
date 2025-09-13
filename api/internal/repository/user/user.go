package user

import (
	"context"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func (userRepo *userRepo) Create(ctx context.Context, user *userDom.User) (*userDom.User, error) {
	userModel, err := NewUserModel(user)
	if err != nil {
		return nil, err
	}

	if err := userRepo.db.WithContext(ctx).Create(userModel).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBCreateError, "failed to create user", err, nil)
	}

	user, err = userModel.ToDomain(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (userRepo *userRepo) Update(ctx context.Context, user *userDom.User) error {
	// Implementation for updating a user
	return nil
}

func (userRepo *userRepo) GetByEmailOrUsername(ctx context.Context, email, username string) (*userDom.User, error) {
	var user User
	if err := userRepo.db.WithContext(ctx).Where("email = ? OR username = ?", email, username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorDom.Raise(ctx, errorDom.ErrUserNotFound, "", err, nil)
		}
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to get user", err, nil)
	}
	domUser, err := user.ToDomain(ctx)
	if err != nil {
		return nil, err
	}
	return domUser, nil
}

func New(db *gorm.DB) userDom.Repository {
	return &userRepo{
		db: db,
	}
}
