package challenge

import (
	"context"

	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/utils"
)

type challengeImpl struct {
	repo challengeDom.Repository
	cvUc cvDom.Usecase
}

func (c *challengeImpl) List(ctx context.Context) ([]*challengeDom.ChallengeDTO, error) {
	domChallenges, err := c.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var challenges []*challengeDom.ChallengeDTO
	for _, challenge := range domChallenges {
		challenges = append(challenges, challenge.ToDTO())
	}

	return challenges, nil
}

func (c *challengeImpl) ValidateAttempt(ctx context.Context, challengeID int64, teamID int64) (*challengeDom.Challenge, error) {
	challenge, err := c.repo.GetByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}

	// check if challenge is visible
	if !challenge.IsVisible {
		return nil, errorDom.Raise(ctx, errorDom.ErrChallengeNotFound, "", nil, nil)
	}

	// fetch submission stats
	subStats, err := c.repo.GetSubmissionStats(ctx, teamID, challengeID)
	if err != nil {
		return nil, err
	}

	// check if already solved
	if subStats.CorrectCount > 0 {
		return nil, errorDom.Raise(ctx, errorDom.ErrChallengeAlreadySolved, "", nil, nil)
	}

	// check if team has exhausted attempts
	if challenge.MaxAttempts > 0 && subStats.IncorrectCount >= challenge.MaxAttempts {
		return nil, errorDom.Raise(ctx, errorDom.ErrChallengeMaxAttemptsReached, "", nil, nil)
	}

	return challenge, nil
}

func (c *challengeImpl) SubmitFlag(ctx context.Context, input *challengeDom.SubmitFlagInput) (*challengeDom.SubmitFlagOutput, error) {
	userID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyUserID)
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	challenge, err := c.ValidateAttempt(ctx, input.ChallengeID, teamID)
	if err != nil {
		return nil, err
	}

	submission := &challengeDom.Submission{
		ChallengeID: input.ChallengeID,
		UserID:      userID,
		TeamID:      teamID,
		Flag:        input.Flag,
		IsCorrect:   input.Flag == challenge.Flag,
		IPAddress:   "",
	}

	if !c.cvUc.GetBool(ctx, cvDom.PostEvent) {
		err = c.repo.SubmitFlag(ctx, submission)
		if err != nil {
			return nil, err
		}
	}

	return &challengeDom.SubmitFlagOutput{
		IsCorrect: input.Flag == challenge.Flag,
	}, nil
}

func New(repo challengeDom.Repository, cvUc cvDom.Usecase) challengeDom.Usecase {
	return &challengeImpl{
		repo: repo,
		cvUc: cvUc,
	}
}
