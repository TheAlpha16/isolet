package user

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/api/infra/database/postgres"
	"github.com/TheAlpha16/isolet/api/internal/domain"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type User struct {
	postgres.BaseModel
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	TeamID   *int64 `gorm:"column:team_id"`
	Role     string `gorm:"type:user_role;default:'player';not null"`
	IsBanned bool   `gorm:"default:false;not null"`
}

func (u *User) ToDomain(ctx context.Context) (*userDom.User, error) {
	role := userDom.Role(u.Role)
	if !role.IsValid() {
		return nil, errorDom.Raise(ctx, errorDom.ErrUserInvalidRole, "", nil, common.ExtraData{
			"role": u.Role,
		})
	}

	return &userDom.User{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(u.CreatedAt, 0),
			UpdatedAt: time.Unix(u.UpdatedAt, 0),
		},
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		TeamID:   u.TeamID,
		Role:     role,
		IsBanned: u.IsBanned,
	}, nil
}

func NewUserModel(u *userDom.User) (*User, error) {
	user := &User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		TeamID:   u.TeamID,
		Role:     string(u.Role),
		IsBanned: u.IsBanned,
	}

	if u.ID != 0 {
		user.ID = u.ID
	}

	return user, nil
}
