package repository

import (
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	cvRepo "github.com/TheAlpha16/isolet/api/internal/repository/configvars"
	userRepo "github.com/TheAlpha16/isolet/api/internal/repository/user"

	"gorm.io/gorm"
)

type Repositories struct {
	User       userDom.Repository
	ConfigVars cvDom.Repository
}

func New(db *gorm.DB) *Repositories {
	userRepo := userRepo.New(db)
	cvRepo := cvRepo.New(db)

	return &Repositories{
		User:       userRepo,
		ConfigVars: cvRepo,
	}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&userRepo.User{}, &cvRepo.ConfigVars{})
}
