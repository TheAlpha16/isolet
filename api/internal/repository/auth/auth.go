package auth

import (
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"

	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) authDom.Repository {
	return &AuthRepo{
		db: db,
	}
}
