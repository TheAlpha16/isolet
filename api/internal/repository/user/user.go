package user

import (
	"context"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"

	"github.com/jackc/pgx/v5/pgconn"
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

	return userModel.ToDomain(ctx)
}

func (userRepo *userRepo) Update(ctx context.Context, user *userDom.User, fields []string) error {
	userModel, err := NewUserModel(user)
	if err != nil {
		return err
	}

	if err := userRepo.db.WithContext(ctx).Model(userModel).Where("id = ?", user.ID).Select(fields).Updates(userModel).Error; err != nil {
		return errorDom.Raise(ctx, errorDom.ErrDBUpdateError, "failed to update user", err, nil)
	}

	return nil
}

func (userRepo *userRepo) GetByEmailOrUsername(ctx context.Context, email, username string) (*userDom.User, error) {
	var user User
	db := userRepo.db.WithContext(ctx)

	if email != "" && username != "" {
		db = db.Where("email = ? OR username = ?", email, username)
	} else if email != "" {
		db = db.Where("email = ?", email)
	} else if username != "" {
		db = db.Where("username = ?", username)
	}

	if err := db.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorDom.Raise(ctx, errorDom.ErrUserNotFound, "", err, nil)
		}
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to get user", err, nil)
	}
	return user.ToDomain(ctx)
}

func (userRepo *userRepo) GetByID(ctx context.Context, id int64) (*userDom.User, error) {
	var user User
	if err := userRepo.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorDom.Raise(ctx, errorDom.ErrUserNotFound, "", err, nil)
		}
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to get user", err, nil)
	}
	return user.ToDomain(ctx)
}

func (userRepo *userRepo) JoinTeam(ctx context.Context, userID int64, teamID int64, teamLen int) error {
	if err := userRepo.db.WithContext(ctx).Model(&User{}).
		Exec("SELECT join_team(?, ?, ?)", userID, teamID, teamLen).Error; err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Message == string(errorDom.ErrTeamFull) {
				return errorDom.Raise(ctx, errorDom.ErrTeamFull, "", err, nil)
			}
		}
		return errorDom.Raise(ctx, errorDom.ErrDBExecError, "failed to join team", err, nil)
	}
	return nil
}

func New(db *gorm.DB) userDom.Repository {
	return &userRepo{
		db: db,
	}
}
