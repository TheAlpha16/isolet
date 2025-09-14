package repository

import (
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	cvRepo "github.com/TheAlpha16/isolet/api/internal/repository/configvars"
	teamRepo "github.com/TheAlpha16/isolet/api/internal/repository/team"
	userRepo "github.com/TheAlpha16/isolet/api/internal/repository/user"

	"gorm.io/gorm"
)

type Repositories struct {
	User       userDom.Repository
	Team       teamDom.Repository
	ConfigVars cvDom.Repository
}

func New(db *gorm.DB) *Repositories {
	userRepo := userRepo.New(db)
	teamRepo := teamRepo.New(db)
	cvRepo := cvRepo.New(db)

	return &Repositories{
		User:       userRepo,
		Team:       teamRepo,
		ConfigVars: cvRepo,
	}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&userRepo.User{}, &cvRepo.ConfigVars{})
}
