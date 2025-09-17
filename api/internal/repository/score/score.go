package score

import (
	"context"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	scoreDom "github.com/TheAlpha16/isolet/api/internal/domain/score"
	challengeRepo "github.com/TheAlpha16/isolet/api/internal/repository/challenge"

	"gorm.io/gorm"
)

type scoreRepo struct {
	db *gorm.DB
}

func (scoreRepo *scoreRepo) GetSubmissionStats(ctx context.Context, teamID int64, challengeIDs []int64) (map[int64]*scoreDom.SubmissionStats, error) {
	result := make(map[int64]*scoreDom.SubmissionStats, len(challengeIDs))
	var rows []SubmissionStatsDTO

	if len(challengeIDs) == 0 {
		return result, nil
	}

	if err := scoreRepo.db.WithContext(ctx).Raw("SELECT challenge_id, COUNT(*) FILTER (WHERE is_correct = true)  AS correct_count, COUNT(*) FILTER (WHERE is_correct = false) AS incorrect_count FROM submissions WHERE team_id = ? AND challenge_id IN ? GROUP BY challenge_id", teamID, challengeIDs).Scan(&rows).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve submission stats", err, common.ExtraData{"team_id": teamID, "challenge_ids": challengeIDs})
	}

	for _, row := range rows {
		result[row.ChallengeID] = &scoreDom.SubmissionStats{
			CorrectCount:   int(row.CorrectCount),
			IncorrectCount: int(row.IncorrectCount),
		}
	}

	return result, nil
}

func (scoreRepo *scoreRepo) GetChallengeSolveCounts(ctx context.Context, challengeIDs []int64) (map[int64]int, error) {
	result := make(map[int64]int, len(challengeIDs))
	var rows []ChallengeSolveCountDTO

	if len(challengeIDs) == 0 {
		return result, nil
	}

	if err := scoreRepo.db.WithContext(ctx).Raw("SELECT challenge_id, COUNT(*) AS solve_count FROM solves WHERE challenge_id IN ? GROUP BY challenge_id", challengeIDs).Scan(&rows).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve challenge solve counts", err, common.ExtraData{"challenge_ids": challengeIDs})
	}

	for _, row := range rows {
		result[row.ChallengeID] = int(row.SolveCount)
	}

	return result, nil
}

func (scoreRepo *scoreRepo) GetTeamSolves(ctx context.Context, teamID int64) (map[int64]struct{}, error) {
	result := make(map[int64]struct{}, 100)
	var rows []*challengeRepo.Solve

	if err := scoreRepo.db.WithContext(ctx).Select("challenge_id").Where("team_id = ?", teamID).Find(&rows).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve team solves", err, common.ExtraData{"team_id": teamID})
	}

	for _, row := range rows {
		result[row.ChallengeID] = struct{}{}
	}

	return result, nil
}

func New(db *gorm.DB) scoreDom.Repository {
	return &scoreRepo{
		db: db,
	}
}
