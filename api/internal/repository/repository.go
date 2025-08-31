package repository

import (
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	authRepo "github.com/TheAlpha16/isolet/api/internal/repository/auth"
	userRepo "github.com/TheAlpha16/isolet/api/internal/repository/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User userDom.Repository
	Auth authDom.Repository
}

func New(db *pgxpool.Pool) *Repositories {
	userRepo := userRepo.New(db)
	authRepo := authRepo.New(db)

	return &Repositories{
		User: userRepo,
		Auth: authRepo,
	}
}
