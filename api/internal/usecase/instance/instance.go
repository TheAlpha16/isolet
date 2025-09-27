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
	cache       cache.Cache
	challengeUc challengeDom.Usecase
}

func (i *instanceImpl) Start(ctx context.Context, input *instanceDom.StartInput) (*instanceDom.Instance, error) {
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
		return instance, nil
	}

	// cache instance start attempt
	if err := i.cacheInstanceAttempt(ctx, teamID, input.ChallengeID); err != nil {
		return nil, err
	}
	defer i.deleteInstanceAttemptCache(ctx, teamID, input.ChallengeID)

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

	return instance, nil
}

func (i *instanceImpl) cacheInstanceAttempt(ctx context.Context, teamID, challengeID int64) error {
	config := utils.GetConfig()
	key := instanceDom.InstanceCacheKey(teamID, challengeID)

	ok, err := i.cache.SetNXWithTTL(ctx, key, "1", config.Instances.CacheTTL)
	if err != nil {
		return err
	}

	// if the key already exists, instance is already being started
	if !ok {
		return errorDom.Raise(ctx, errorDom.ErrInstanceAlreadyStarting, "instance is already being started", nil, common.ExtraData{"team_id": teamID, "challenge_id": challengeID})
	}

	return nil
}

func (i *instanceImpl) deleteInstanceAttemptCache(ctx context.Context, teamID, challengeID int64) error {
	return i.cache.Delete(ctx, instanceDom.InstanceCacheKey(teamID, challengeID))
}

func New(repo instanceDom.Repository, cache cache.Cache, challengeUc challengeDom.Usecase) instanceDom.Usecase {
	return &instanceImpl{
		repo:        repo,
		cache:       cache,
		challengeUc: challengeUc,
	}
}
