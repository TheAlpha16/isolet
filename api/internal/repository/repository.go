package repository

import (
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"

	userRepo "github.com/TheAlpha16/isolet/api/internal/repository/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User userDom.Repository
}

func New(db *pgxpool.Pool) *Repositories {
	userRepo := userRepo.New(db)

	return &Repositories{
		User: userRepo,
	}
}
