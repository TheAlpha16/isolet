package score

import (
	"context"

	scoreDom "github.com/TheAlpha16/isolet/api/internal/domain/score"
)

type scoreImpl struct {
	repo scoreDom.Repository
}

func (scoreImpl *scoreImpl) GetSubmissionStats(ctx context.Context, teamID int64, challengeIDs []int64) (map[int64]*scoreDom.SubmissionStats, error) {
	return scoreImpl.repo.GetSubmissionStats(ctx, teamID, challengeIDs)
}

func (scoreImpl *scoreImpl) GetChallengeSolveCounts(ctx context.Context, challengeIDs []int64) (map[int64]int, error) {
	return scoreImpl.repo.GetChallengeSolveCounts(ctx, challengeIDs)
}

func (scoreImpl *scoreImpl) GetTeamSolves(ctx context.Context, teamID int64) (map[int64]struct{}, error) {
	return scoreImpl.repo.GetTeamSolves(ctx, teamID)
}

func (scoreImpl *scoreImpl) GetTeamScore(ctx context.Context, teamID int64) (int, error) {
	return scoreImpl.repo.GetTeamScore(ctx, teamID)
}

func New(repo scoreDom.Repository) scoreDom.Usecase {
	return &scoreImpl{
		repo: repo,
	}
}
