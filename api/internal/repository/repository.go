package repository

import (
	domain_user "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/jackc/pgx/v5/pgxpool"

	repo_user "github.com/TheAlpha16/isolet/api/internal/repository/user"
)

type Repositories struct {
	User domain_user.Repository
}

func New(db *pgxpool.Pool) *Repositories {
	userRepo := repo_user.New(db)

	return &Repositories{
		User: userRepo,
	}
}
