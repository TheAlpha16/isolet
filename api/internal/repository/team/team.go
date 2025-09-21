package team

import (
	"context"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	userRepo "github.com/TheAlpha16/isolet/api/internal/repository/user"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type teamRepo struct {
	db *gorm.DB
}

func (teamRepo *teamRepo) Create(ctx context.Context, team *teamDom.Team) (*teamDom.Team, error) {
	teamModel, err := NewTeamModel(team)
	if err != nil {
		return nil, err
	}

	var domainTeam *teamDom.Team
	err = teamRepo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// create the team
		if err := tx.Create(teamModel).Error; err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				if pgErr.Code == errorDom.PgErrDuplicateKey {
					return errorDom.Raise(ctx, errorDom.ErrTeamNameTaken, "", err, nil)
				}
			}
			return errorDom.Raise(ctx, errorDom.ErrDBCreateError, "failed to create team", err, nil)
		}

		// update captain's team_id and role
		if err := tx.Model(&userRepo.User{}).
			Where("id = ?", team.CaptainID).
			Updates(map[string]any{
				"team_id": teamModel.ID,
				"role":    userDom.RoleCaptain,
			}).Error; err != nil {
			return errorDom.Raise(ctx, errorDom.ErrDBUpdateError, "failed to assign captain to team", err, nil)
		}

		domainTeam, err = teamModel.ToDomain(ctx)
		return err
	})

	if err != nil {
		return nil, err
	}

	return domainTeam, nil
}

func (teamRepo *teamRepo) GetByName(ctx context.Context, name string) (*teamDom.Team, error) {
	var team Team
	if err := teamRepo.db.WithContext(ctx).Where("name = ?", name).First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorDom.Raise(ctx, errorDom.ErrTeamNotFound, "", err, nil)
		}
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to get team by name", err, nil)
	}

	domainTeam, err := team.ToDomain(ctx)
	if err != nil {
		return nil, err
	}

	return domainTeam, nil
}

func (teamRepo *teamRepo) GetByID(ctx context.Context, id int64) (*teamDom.Team, error) {
	var team Team
	if err := teamRepo.db.WithContext(ctx).Preload("Members").First(&team, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorDom.Raise(ctx, errorDom.ErrTeamNotFound, "", err, nil)
		}
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to get team by id", err, nil)
	}

	return team.ToDomain(ctx)
}

func (teamRepo *teamRepo) GetNameByIDs(ctx context.Context, teamIDs []int64) (map[int64]string, error) {
	var teams []Team

	if len(teamIDs) == 0 {
		return make(map[int64]string), nil
	}

	if err := teamRepo.db.WithContext(ctx).Select("id, name").Where("id IN ?", teamIDs).Find(&teams).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to get team names by ids", err, nil)
	}

	names := make(map[int64]string)
	for _, team := range teams {
		names[team.ID] = team.Name
	}

	return names, nil
}

func New(db *gorm.DB) teamDom.Repository {
	return &teamRepo{
		db: db,
	}
}
