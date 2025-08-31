package auth

import (
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) authDom.Repository {
	return &AuthRepo{
		db: db,
	}
}
