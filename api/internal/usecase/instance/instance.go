package instance

import (
	"context"

	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	"github.com/TheAlpha16/isolet/api/utils"
)

type instanceImpl struct {
	repo        instanceDom.Repository
	challengeUc challengeDom.Usecase
}

func (i *instanceImpl) Start(ctx context.Context, input *instanceDom.StartInput) (*instanceDom.Instance, error) {
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	challenge, err := i.challengeUc.ValidateAccess(ctx, input.ChallengeID, teamID)
	if err != nil {
		return nil, err
	}
	if challenge.Type != challengeDom.ChallengeOnDemand {
		return nil, errorDom.Raise(ctx, errorDom.ErrInstanceNotOnDemand, "", nil, nil)
	}
	return nil, nil
}

func New(repo instanceDom.Repository, challengeUc challengeDom.Usecase) instanceDom.Usecase {
	return &instanceImpl{
		repo:        repo,
		challengeUc: challengeUc,
	}
}
