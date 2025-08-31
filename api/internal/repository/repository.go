package repository

import (
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	authRepo "github.com/TheAlpha16/isolet/api/internal/repository/auth"
	userRepo "github.com/TheAlpha16/isolet/api/internal/repository/user"

	"gorm.io/gorm"
)

type Repositories struct {
	User userDom.Repository
	Auth authDom.Repository
}

func New(db *gorm.DB) *Repositories {
	userRepo := userRepo.New(db)
	authRepo := authRepo.New(db)

	return &Repositories{
		User: userRepo,
		Auth: authRepo,
	}
}
