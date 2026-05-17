package repository

import (
	"github.com/TheAlpha16/isolet/oracle/infra/cache"
	challengeDom "github.com/TheAlpha16/isolet/oracle/internal/domain/challenge"
	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/oracle/internal/domain/manifest"
	scoreDom "github.com/TheAlpha16/isolet/oracle/internal/domain/score"
	teamDom "github.com/TheAlpha16/isolet/oracle/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/oracle/internal/domain/user"
	challengeRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/challenge"
	cvRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/configvars"
	instanceRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/instance"
	manifestRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/manifest"
	scoreRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/score"
	teamRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/team"
	userRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/user"

	"gorm.io/gorm"
)

type Repositories struct {
	User       userDom.Repository
	Team       teamDom.Repository
	ConfigVars cvDom.Repository
	Challenge  challengeDom.Repository
	Score      scoreDom.Repository
	Instance   instanceDom.Repository
	Manifest   manifestDom.Repository
}

func New(db *gorm.DB, cacheClient cache.Cache) *Repositories {
	users := userRepo.New(db)
	teams := teamRepo.New(db)
	cv := cvRepo.New(db)
	challenges := challengeRepo.New(db, cacheClient)
	scores := scoreRepo.New(db)
	instances := instanceRepo.New(db)
	manifests := manifestRepo.New(db, cacheClient)

	return &Repositories{
		User:       users,
		Team:       teams,
		ConfigVars: cv,
		Challenge:  challenges,
		Score:      scores,
		Instance:   instances,
		Manifest:   manifests,
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
		&challengeRepo.UnlockedHint{},
		&instanceRepo.Instance{},
		&instanceRepo.Endpoint{},
		&manifestRepo.Manifest{},
		&manifestRepo.Resource{},
		&manifestRepo.EndpointSpec{},
	)
}
