package user

import (
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func (userRepo *UserRepo) Create(user *userDom.User) error {
	// Implementation for creating a user
	return nil
}

func (userRepo *UserRepo) Update(user *userDom.User) error {
	// Implementation for updating a user
	return nil
}

func New(db *pgxpool.Pool) userDom.Repository {
	return &UserRepo{
		db: db,
	}
}
