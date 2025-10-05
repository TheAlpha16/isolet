package instance

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/api/infra/cache"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	"github.com/TheAlpha16/isolet/api/utils"
)

type instanceImpl struct {
	repo        instanceDom.Repository
	service     instanceDom.Service
	cache       cache.Cache
	challengeUc challengeDom.Usecase
}

func (i *instanceImpl) Start(ctx context.Context, input *instanceDom.StartInput) (*instanceDom.InstanceDTO, error) {
	config := utils.GetConfig()
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	// check if the challenge exists for this team and is of type on-demand
	challenge, err := i.challengeUc.ValidateAccess(ctx, input.ChallengeID, teamID)
	if err != nil {
		return nil, err
	}
	if challenge.Type != challengeDom.ChallengeOnDemand {
		return nil, errorDom.Raise(ctx, errorDom.ErrInstanceNotOnDemand, "", nil, nil)
	}

	// check if the instance is already running for the team and challenge
	instance, err := i.repo.GetByTeamAndChallenge(ctx, teamID, input.ChallengeID)
	if err != nil {
		if !errorDom.IsSameError(err, errorDom.ErrInstanceNotFound) {
			return nil, err
		}
	}
	if instance != nil {
		return instance.ToDTO(), nil
	}

	// take mutex on the instance to prevent race conditions
	if err := i.acquireInstanceLock(ctx, teamID, input.ChallengeID, config.Instances.StartTimeout); err != nil {
		return nil, err
	}
	defer i.releaseInstanceLock(ctx, teamID, input.ChallengeID)

	// TODO create instance

	// insert into database
	instance, err = i.repo.Create(ctx, &instanceDom.Instance{
		TeamID:      teamID,
		ChallengeID: input.ChallengeID,
		ExpiresAt:   time.Now().Add(utils.GetConfig().Instances.Lifetime).Unix(),
	})
	if err != nil {
		return nil, err
	}

	return instance.ToDTO(), nil
}

func (i *instanceImpl) Stop(ctx context.Context, input *instanceDom.StopInput) error {
	config := utils.GetConfig()
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	// check if the instance exists for the team and challenge
	instance, err := i.repo.GetByID(ctx, input.InstanceID)
	if err != nil {
		return err
	}
	if instance.TeamID != teamID {
		return errorDom.Raise(ctx, errorDom.ErrInstanceNotFound, "", nil, nil)
	}

	// take mutex on the instance to prevent race conditions
	if err := i.acquireInstanceLock(ctx, teamID, instance.ChallengeID, config.Instances.StopTimeout); err != nil {
		return err
	}
	defer i.releaseInstanceLock(ctx, teamID, instance.ChallengeID)

	// TODO stop instance

	// delete from database
	if err := i.repo.Delete(ctx, instance.ID); err != nil {
		return err
	}

	return nil
}

func (i *instanceImpl) Extend(ctx context.Context, input *instanceDom.ExtendInput) (*instanceDom.InstanceDTO, error) {
	config := utils.GetConfig()
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	// check if the instance exists for the team and challenge
	instance, err := i.repo.GetByID(ctx, input.InstanceID)
	if err != nil {
		return nil, err
	}
	if instance.TeamID != teamID {
		return nil, errorDom.Raise(ctx, errorDom.ErrInstanceNotFound, "", nil, nil)
	}

	// take mutex on the instance to prevent race conditions
	if err := i.acquireInstanceLock(ctx, teamID, instance.ChallengeID, config.Instances.ExtendTimeout); err != nil {
		return nil, err
	}
	defer i.releaseInstanceLock(ctx, teamID, instance.ChallengeID)

	// TODO extend instance deadline

	// update in database
	instance.ExpiresAt = time.Unix(instance.ExpiresAt, 0).Add(config.Instances.Lifetime).Unix()
	if err := i.repo.Update(ctx, instance, []string{"expires_at"}); err != nil {
		return nil, err
	}

	return instance.ToDTO(), nil
}

func (i *instanceImpl) acquireInstanceLock(ctx context.Context, teamID, challengeID int64, ttl time.Duration) error {
	key := instanceDom.InstanceCacheKey(teamID, challengeID)

	ok, err := i.cache.SetNXWithTTL(ctx, key, "1", ttl)
	if err != nil {
		return err
	}

	// if the key already exists, instance is already being started/stopped/updated
	if !ok {
		return errorDom.Raise(ctx, errorDom.ErrInstanceWIP, "", nil, common.ExtraData{"team_id": teamID, "challenge_id": challengeID})
	}

	return nil
}

func (i *instanceImpl) releaseInstanceLock(ctx context.Context, teamID, challengeID int64) error {
	return i.cache.Delete(ctx, instanceDom.InstanceCacheKey(teamID, challengeID))
}

func New(repo instanceDom.Repository, service instanceDom.Service, cache cache.Cache, challengeUc challengeDom.Usecase) instanceDom.Usecase {
	return &instanceImpl{
		repo:        repo,
		service:     service,
		cache:       cache,
		challengeUc: challengeUc,
	}
}
