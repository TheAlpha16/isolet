package challenge

import (
	"context"

	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	scoreDom "github.com/TheAlpha16/isolet/api/internal/domain/score"
	"github.com/TheAlpha16/isolet/api/utils"
)

type challengeImpl struct {
	repo    challengeDom.Repository
	cvUc    cvDom.Usecase
	scoreUc scoreDom.Usecase
}

func (c *challengeImpl) List(ctx context.Context) ([]*challengeDom.ChallengeDTO, error) {
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)
	filteredChallenges := make([]*challengeDom.Challenge, 0)

	domChallenges, err := c.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	solves, err := c.scoreUc.GetTeamSolves(ctx, teamID)
	if err != nil {
		return nil, err
	}

	unlockedHints, err := c.repo.GetUnlockedHints(ctx, teamID)
	if err != nil {
		return nil, err
	}

	var challengeIDs []int64
	for _, challenge := range domChallenges {
		if !challenge.AreRequirementsMet(solves) || !challenge.IsVisible || !challenge.Category.IsVisible {
			continue
		}
		challengeIDs = append(challengeIDs, challenge.ID)
		filteredChallenges = append(filteredChallenges, challenge)

		for _, hint := range challenge.Hints {
			if hint.IsUnlocked(unlockedHints) {
				hint.Unlocked = true
				continue
			}
			hint.Text = ""
		}
	}

	challengeSolveCounts, err := c.scoreUc.GetChallengeSolveCounts(ctx, challengeIDs)
	if err != nil {
		return nil, err
	}

	subStats, err := c.scoreUc.GetSubmissionStats(ctx, teamID, challengeIDs)
	if err != nil {
		return nil, err
	}

	challenges := make([]*challengeDom.ChallengeDTO, 0, len(filteredChallenges))
	for _, challenge := range filteredChallenges {
		challengeDTO := challenge.ToDTO()
		enrichChallengeDTO(challengeDTO, subStats, challengeSolveCounts)
		challenges = append(challenges, challengeDTO)
	}

	return challenges, nil
}

func (c *challengeImpl) ValidateAttempt(ctx context.Context, challengeID int64, teamID int64) (*challengeDom.Challenge, error) {
	challenge, err := c.repo.GetByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}

	// check if challenge is visible
	if !challenge.IsVisible || !challenge.Category.IsVisible {
		return nil, errorDom.Raise(ctx, errorDom.ErrChallengeNotFound, "", nil, nil)
	}

	// fetch submission stats
	subStats, err := c.scoreUc.GetSubmissionStats(ctx, teamID, []int64{challengeID})
	if err != nil {
		return nil, err
	}
	challengeStats, ok := subStats[challengeID]
	if !ok {
		challengeStats = &scoreDom.SubmissionStats{}
	}

	// check if already solved
	if challengeStats.CorrectCount > 0 {
		return nil, errorDom.Raise(ctx, errorDom.ErrChallengeAlreadySolved, "", nil, nil)
	}

	// check if team has exhausted attempts
	if challenge.MaxAttempts > 0 && challengeStats.IncorrectCount >= challenge.MaxAttempts {
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
		Points:      challenge.Points,
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

func (c *challengeImpl) UnlockHint(ctx context.Context, input *challengeDom.UnlockHintInput) (*challengeDom.Hint, error) {
	/*
		 	* 1. if already unlocked
			* 2. does hint exist
			* 3. is hint visible
			* 4. is challenge visible
			* 5. are challenge requirements met
			* 6. does team have enough score
			*
		* fetch hint
		* fetch linked challenge
		* fetch team solves
		* fetch team score
		* insert unlocked hint row
	*/
	return nil, nil
}

func enrichChallengeDTO(challengeDTO *challengeDom.ChallengeDTO, submissionStatsMap map[int64]*scoreDom.SubmissionStats, challengeSolveCounts map[int64]int) {
	challengeDTO.TotalSolves = challengeSolveCounts[challengeDTO.ID]

	stats := submissionStatsMap[challengeDTO.ID]
	if stats == nil {
		return
	}

	challengeDTO.Solved = stats.CorrectCount > 0
	challengeDTO.AttemptCount = stats.IncorrectCount + stats.CorrectCount
}

func New(repo challengeDom.Repository, cvUc cvDom.Usecase, scoreUc scoreDom.Usecase) challengeDom.Usecase {
	return &challengeImpl{
		repo:    repo,
		cvUc:    cvUc,
		scoreUc: scoreUc,
	}
}
