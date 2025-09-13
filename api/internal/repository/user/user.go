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

func (userRepo *userRepo) CheckIdentifiers(ctx context.Context, email, username string) (*userDom.IdentifierExistence, error) {
	var exists userDom.IdentifierExistence
	query := `SELECT BOOL_OR(email = $1) AS email_exists, BOOL_OR(username = $2) AS username_exists FROM users WHERE email = $1 OR username = $2;`

	if err := userRepo.db.WithContext(ctx).Raw(query, email, username).Scan(&exists).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to check user identifiers", err, nil)
	}
	return &exists, nil
}

func New(db *gorm.DB) userDom.Repository {
	return &userRepo{
		db: db,
	}
}
