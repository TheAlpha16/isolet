package repository

import (
	"github.com/TheAlpha16/isolet/api/infra/cache"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	scoreDom "github.com/TheAlpha16/isolet/api/internal/domain/score"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	challengeRepo "github.com/TheAlpha16/isolet/api/internal/repository/challenge"
	cvRepo "github.com/TheAlpha16/isolet/api/internal/repository/configvars"
	scoreRepo "github.com/TheAlpha16/isolet/api/internal/repository/score"
	teamRepo "github.com/TheAlpha16/isolet/api/internal/repository/team"
	userRepo "github.com/TheAlpha16/isolet/api/internal/repository/user"

	"gorm.io/gorm"
)

type Repositories struct {
	User       userDom.Repository
	Team       teamDom.Repository
	ConfigVars cvDom.Repository
	Challenge  challengeDom.Repository
	Score      scoreDom.Repository
}

func New(db *gorm.DB, cache cache.Cache) *Repositories {
	userRepo := userRepo.New(db)
	teamRepo := teamRepo.New(db)
	cvRepo := cvRepo.New(db)
	challengeRepo := challengeRepo.New(db, cache)
	scoreRepo := scoreRepo.New(db)

	return &Repositories{
		User:       userRepo,
		Team:       teamRepo,
		ConfigVars: cvRepo,
		Challenge:  challengeRepo,
		Score:      scoreRepo,
	}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&cvRepo.ConfigVars{},
		&userRepo.User{},
		&teamRepo.Team{},
		&challengeRepo.Category{},
		&challengeRepo.Challenge{},
		&challengeRepo.Hint{},
		&challengeRepo.Submission{},
		&challengeRepo.Solve{},
	)
}
