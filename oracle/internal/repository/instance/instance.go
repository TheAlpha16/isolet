package instance

import (
	"context"
	"errors"
	"time"

	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	"github.com/TheAlpha16/isolet/oracle/utils"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type instanceRepo struct {
	db *gorm.DB
}

func (ir *instanceRepo) Create(ctx context.Context, instance *instanceDom.Instance) (*instanceDom.Instance, error) {
	extraData := common.ExtraData{utils.ContextKeyTeamID: instance.TeamID, utils.ContextKeyChallengeID: instance.ChallengeID}
	instanceModel, err := NewInstanceModel(instance)
	if err != nil {
		return nil, err
	}

	if err := ir.db.WithContext(ctx).Create(instanceModel).Error; err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == errorDom.PgErrDuplicateKey {
				return nil, errorDom.Raise(ctx, errorDom.ErrInstanceAlreadyRunning, "", err, extraData)
			}
		}
		return nil, errorDom.Raise(ctx, errorDom.ErrDBCreateError, "failed to create instance", err, extraData)
	}

	return instanceModel.ToDomain(ctx)
}

func (ir *instanceRepo) GetByID(ctx context.Context, id int64) (*instanceDom.Instance, error) {
	var instance Instance
	if err := ir.db.WithContext(ctx).Preload("Endpoints").First(&instance, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorDom.Raise(ctx, errorDom.ErrInstanceNotFound, "", err, common.ExtraData{"id": id})
		}
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to get instance by id", err, common.ExtraData{"id": id})
	}

	return instance.ToDomain(ctx)
}

func (ir *instanceRepo) Update(ctx context.Context, id int64, updates map[string]any) error {
	if err := ir.db.WithContext(ctx).Model(&Instance{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return errorDom.Raise(ctx, errorDom.ErrDBUpdateError, "failed to update instance", err, common.ExtraData{"id": id})
	}

	return nil
}

func (ir instanceRepo) Delete(ctx context.Context, id int64) error {
	if err := ir.db.WithContext(ctx).Delete(&Instance{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorDom.Raise(ctx, errorDom.ErrInstanceNotFound, "", err, nil)
		}
		return errorDom.Raise(ctx, errorDom.ErrDBDeleteError, "failed to delete instance", err, common.ExtraData{"id": id})
	}
	return nil
}

func (ir *instanceRepo) GetByTeam(ctx context.Context, teamID int64) ([]*instanceDom.Instance, error) {
	var instances []Instance
	if err := ir.db.WithContext(ctx).Preload("Endpoints").Where("team_id = ?", teamID).Find(&instances).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to get instances by team", err, common.ExtraData{utils.ContextKeyTeamID: teamID})
	}

	domainInstances := make([]*instanceDom.Instance, 0, len(instances))
	for _, inst := range instances {
		domainInst, err := inst.ToDomain(ctx)
		if err != nil {
			return nil, err
		}
		domainInstances = append(domainInstances, domainInst)
	}
	return domainInstances, nil
}

func (ir *instanceRepo) GetByRefs(ctx context.Context, teamID *int64, challengeID int64) (*instanceDom.Instance, error) {
	var instance Instance
	query := ir.db.WithContext(ctx).Preload("Endpoints").Where("challenge_id = ?", challengeID)

	if teamID != nil {
		query = query.Where("team_id = ?", *teamID)
	} else {
		query = query.Where("team_id IS NULL")
	}

	if err := query.First(&instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorDom.Raise(ctx, errorDom.ErrInstanceNotFound, "", err, nil)
		}
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to get instance by refs", err, common.ExtraData{utils.ContextKeyTeamID: teamID, utils.ContextKeyChallengeID: challengeID})
	}

	return instance.ToDomain(ctx)
}

func (ir *instanceRepo) DeleteByRefs(ctx context.Context, teamID *int64, challengeID int64) (int64, error) {
	query := ir.db.WithContext(ctx).Where("challenge_id = ?", challengeID)

	if teamID != nil {
		query = query.Where("team_id = ?", *teamID)
	} else {
		query = query.Where("team_id IS NULL")
	}

	resp := query.Delete(&Instance{})
	if err := resp.Error; err != nil {
		return 0, errorDom.Raise(ctx, errorDom.ErrDBDeleteError, "failed to delete instance by refs", err, common.ExtraData{utils.ContextKeyTeamID: teamID, utils.ContextKeyChallengeID: challengeID})
	}
	return resp.RowsAffected, nil
}

func (ir *instanceRepo) DeleteExpired(ctx context.Context, now time.Time) (map[int64]int64, error) {
	var instances []Instance

	if err := ir.db.WithContext(ctx).Clauses(
		clause.Returning{Columns: []clause.Column{{Name: "challenge_id"}}},
	).Where("expires_at <= ?", now.Unix()).Delete(&instances).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBDeleteError, "failed to delete expired instances", err, common.ExtraData{"now": now})
	}

	result := make(map[int64]int64, len(instances))
	for _, inst := range instances {
		result[inst.ChallengeID]++
	}

	return result, nil
}

func New(db *gorm.DB) instanceDom.Repository {
	return &instanceRepo{
		db: db,
	}
}
