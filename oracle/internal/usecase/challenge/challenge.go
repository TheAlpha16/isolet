package challenge

import (
	"context"
	"sync"

	challengeDom "github.com/TheAlpha16/isolet/oracle/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	scoreDom "github.com/TheAlpha16/isolet/oracle/internal/domain/score"
	"github.com/TheAlpha16/isolet/oracle/utils"
)

type challengeImpl struct {
	repo         challengeDom.Repository
	instanceRepo instanceDom.Repository
	cvUc         cvDom.Usecase
	scoreUc      scoreDom.Usecase
	wg           *sync.WaitGroup
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

func (c *challengeImpl) ValidateAccess(ctx context.Context, challengeID, teamID int64) (*challengeDom.Challenge, error) {
	challenge, err := c.repo.GetByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}

	// check if challenge is visible
	if !challenge.IsVisible || !challenge.Category.IsVisible {
		return nil, errorDom.Raise(ctx, errorDom.ErrChallengeNotFound, "", nil, nil)
	}

	// check if requirements are met
	teamSolves, err := c.scoreUc.GetTeamSolves(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if !challenge.AreRequirementsMet(teamSolves) {
		return nil, errorDom.Raise(ctx, errorDom.ErrChallengeNotFound, "", nil, nil)
	}

	return challenge, nil
}

func (c *challengeImpl) SubmitFlag(ctx context.Context, input *challengeDom.SubmitFlagInput) (*challengeDom.SubmitFlagOutput, error) {
	userID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyUserID)
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)
	ipAddress := common.GetFieldFromExtraData[string](ctx, utils.ContextKeyIP)

	challenge, err := c.validateAttempt(ctx, input.ChallengeID, teamID)
	if err != nil {
		return nil, err
	}

	isCorrect, err := c.validateFlag(ctx, challenge, input.Flag)
	if err != nil {
		return nil, err
	}

	submission := &challengeDom.Submission{
		ChallengeID: input.ChallengeID,
		UserID:      userID,
		TeamID:      teamID,
		Flag:        input.Flag,
		IsCorrect:   isCorrect,
		Points:      challenge.Points,
		IPAddress:   ipAddress,
	}

	if !c.cvUc.GetBool(ctx, cvDom.EventPostMode) {
		if err := c.repo.SubmitFlag(ctx, submission); err != nil {
			return nil, err
		}

		if submission.IsCorrect {
			rCtx := context.WithoutCancel(ctx)
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				c.updateScoreboard(rCtx, teamID, submission.Points)
			}()
		}
	}

	return &challengeDom.SubmitFlagOutput{
		IsCorrect: isCorrect,
	}, nil
}

func (c *challengeImpl) UnlockHint(ctx context.Context, input *challengeDom.UnlockHintInput) (*challengeDom.HintDTO, error) {
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)
	ErrHintNotFound := errorDom.Raise(ctx, errorDom.ErrHintNotFound, "", nil, nil)

	hint, err := c.repo.GetHintByID(ctx, input.HintID)
	if err != nil {
		return nil, err
	}
	if !hint.IsVisible {
		return nil, ErrHintNotFound
	}

	challenge, err := c.repo.GetByID(ctx, hint.ChallengeID)
	if err != nil {
		if errorDom.IsSameError(err, errorDom.ErrChallengeNotFound) {
			return nil, ErrHintNotFound
		}
		return nil, err
	}
	if !challenge.IsVisible {
		return nil, ErrHintNotFound
	}

	teamSolves, err := c.scoreUc.GetTeamSolves(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if !challenge.AreRequirementsMet(teamSolves) {
		return nil, ErrHintNotFound
	}

	teamScore, err := c.scoreUc.GetTeamScore(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if teamScore < hint.Cost {
		return nil, errorDom.Raise(ctx, errorDom.ErrHintCostExceeded, "", nil, nil)
	}

	uHint := &challengeDom.UnlockedHint{
		TeamID: teamID,
		HintID: hint.ID,
		Cost:   hint.Cost,
	}
	hint.Unlocked = true

	if err := c.repo.UnlockHint(ctx, uHint); err != nil {
		return nil, err
	}

	// update the scoreboard
	rCtx := context.WithoutCancel(ctx)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.updateScoreboard(rCtx, teamID, -uHint.Cost)
	}()

	return hint.ToDTO(), nil
}

func (c *challengeImpl) validateAttempt(ctx context.Context, challengeID int64, teamID int64) (*challengeDom.Challenge, error) {
	challenge, err := c.ValidateAccess(ctx, challengeID, teamID)
	if err != nil {
		return nil, err
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

func (c *challengeImpl) validateFlag(ctx context.Context, challenge *challengeDom.Challenge, flag string) (bool, error) {
	if challenge.Type == challengeDom.ChallengeOnDemand {
		// fetch the instance for the team and challenge
		teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)
		instance, err := c.instanceRepo.GetByRefs(ctx, &teamID, challenge.ID)
		if err != nil {
			if errorDom.IsSameError(err, errorDom.ErrInstanceNotFound) {
				return false, errorDom.Raise(ctx, errorDom.ErrInstanceNotRunning, "", err, nil)
			}
			return false, err
		}
		if instance != nil && instance.Flag != nil {
			return *instance.Flag == flag, nil
		}
	}
	return flag == challenge.Flag, nil
}

func (c *challengeImpl) GetTeamSubmissions(ctx context.Context, teamID int64) ([]*challengeDom.Submission, error) {
	return c.repo.GetTeamSubmissions(ctx, teamID)
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

func (c *challengeImpl) updateScoreboard(ctx context.Context, teamID int64, delta int) {
	err := c.scoreUc.UpdateTeamScore(ctx, teamID, delta)
	if err == nil {
		return
	}
	errorDom.RaiseToSentry(ctx, err)

	// trigger the refresh scoreboard
	if err := c.scoreUc.RefreshScoreboard(ctx); err != nil {
		errorDom.RaiseToSentry(ctx, err)
	}
}

func New(repo challengeDom.Repository, instanceRepo instanceDom.Repository, cvUc cvDom.Usecase, scoreUc scoreDom.Usecase, wg *sync.WaitGroup) challengeDom.Usecase {
	return &challengeImpl{
		repo:         repo,
		instanceRepo: instanceRepo,
		cvUc:         cvUc,
		scoreUc:      scoreUc,
		wg:           wg,
	}
}
