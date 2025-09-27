package instance

import (
	"context"
	"errors"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type instanceRepo struct {
	db *gorm.DB
}

func (ir *instanceRepo) Create(ctx context.Context, instance *instanceDom.Instance) (*instanceDom.Instance, error) {
	extraData := common.ExtraData{"team_id": instance.TeamID, "challenge_id": instance.ChallengeID}
	instanceModel, err := NewInstanceModel(instance)
	if err != nil {
		return nil, err
	}

	if err := ir.db.WithContext(ctx).Create(instanceModel).Error; err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errorDom.PgErrDuplicateKey {
				return nil, errorDom.Raise(ctx, errorDom.ErrInstanceAlreadyStarting, "", err, extraData)
			}
		}
		return nil, errorDom.Raise(ctx, errorDom.ErrDBCreateError, "failed to create instance", err, extraData)
	}

	return instanceModel.ToDomain(ctx)
}

func (ir *instanceRepo) GetByTeamAndChallenge(ctx context.Context, teamID, challengeID int64) (*instanceDom.Instance, error) {
	var instance Instance
	if err := ir.db.WithContext(ctx).Where("team_id = ? AND challenge_id = ?", teamID, challengeID).First(&instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorDom.Raise(ctx, errorDom.ErrInstanceNotFound, "", err, nil)
		}
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to get instance by team and challenge", err, common.ExtraData{"team_id": teamID, "challenge_id": challengeID})
	}

	return instance.ToDomain(ctx)
}

func New(db *gorm.DB) instanceDom.Repository {
	return &instanceRepo{
		db: db,
	}
}
