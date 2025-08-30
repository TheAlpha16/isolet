package user

import (
	domain_user "github.com/TheAlpha16/isolet/api/internal/domain/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func (userRepo *UserRepo) Create(user *domain_user.User) error {
	// Implementation for creating a user
	return nil
}

func (userRepo *UserRepo) Update(user *domain_user.User) error {
	// Implementation for updating a user
	return nil
}

func New(db *pgxpool.Pool) domain_user.Repository {
	return &UserRepo{
		db: db,
	}
}
