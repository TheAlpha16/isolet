package team

import (
	"context"
	"errors"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	userRepo "github.com/TheAlpha16/isolet/api/internal/repository/user"

	"gorm.io/gorm"
)

type teamRepo struct {
	db *gorm.DB
}

func (t *teamRepo) Create(ctx context.Context, team *teamDom.Team) (*teamDom.Team, error) {
	teamModel, err := NewTeamModel(team)
	if err != nil {
		return nil, err
	}

	var domainTeam *teamDom.Team
	err = t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// create the team
		if err := tx.Create(teamModel).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errorDom.Raise(ctx, errorDom.ErrTeamNameTaken, "", err, nil)
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

func New(db *gorm.DB) teamDom.Repository {
	return &teamRepo{
		db: db,
	}
}
