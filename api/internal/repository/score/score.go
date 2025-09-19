package score

import (
	"context"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	scoreDom "github.com/TheAlpha16/isolet/api/internal/domain/score"
	challengeRepo "github.com/TheAlpha16/isolet/api/internal/repository/challenge"

	"github.com/lib/pq"
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

func (scoreRepo *scoreRepo) GetTeamScore(ctx context.Context, teamID int64) (int, error) {
	var score int

	if err := scoreRepo.db.WithContext(ctx).Raw("SELECT COALESCE((SELECT SUM(points) FROM solves WHERE team_id = ?), 0) - COALESCE((SELECT SUM(cost) FROM unlocked_hints WHERE team_id = ?), 0) AS score;", teamID, teamID).Scan(&score).Error; err != nil {
		return 0, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve team score", err, common.ExtraData{"team_id": teamID})
	}

	return score, nil
}

func (scoreRepo *scoreRepo) GetScoreboard(ctx context.Context) (map[int64]int, error) {
	result := make(map[int64]int, 100)
	var rows []*ScoreboardRow

	if err := scoreRepo.db.WithContext(ctx).Raw("SELECT s.team_id, COALESCE(SUM(s.points), 0) - COALESCE(MAX(h.total_cost), 0) AS score FROM solves s LEFT JOIN ( SELECT team_id, SUM(cost) AS total_cost FROM unlocked_hints GROUP BY team_id ) h ON s.team_id = h.team_id GROUP BY s.team_id;").Scan(&rows).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve scoreboard", err, nil)
	}

	for _, row := range rows {
		result[row.TeamID] = int(row.Score)
	}

	return result, nil
}

func (scoreRepo *scoreRepo) GetTeamsScoreRecords(ctx context.Context, teamIDs []int64) (map[int64][]*scoreDom.ScoreRecord, error) {
	result := make(map[int64][]*scoreDom.ScoreRecord, len(teamIDs))
	var rows []*ScoreRecordRow

	if len(teamIDs) == 0 {
		return result, nil
	}

	if err := scoreRepo.db.WithContext(ctx).Raw("SELECT s.team_id, s.points AS points, s.created_at AS timestamp FROM solves s WHERE s.team_id = ANY(?) UNION ALL SELECT uh.team_id, -uh.cost AS points, uh.created_at AS timestamp FROM unlocked_hints uh WHERE uh.team_id = ANY(?) ORDER BY timestamp;", pq.Array(teamIDs), pq.Array(teamIDs)).Scan(&rows).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve solve records", err, nil)
	}

	for _, row := range rows {
		result[row.TeamID] = append(result[row.TeamID], row.ToDomain(ctx))
	}

	return result, nil
}

func (scoreRepo *scoreRepo) GetTeamScoreRecords(ctx context.Context, teamID int64) ([]*scoreDom.ScoreRecord, error) {
	var result []*scoreDom.ScoreRecord
	var rows []*ScoreRecordRow

	if err := scoreRepo.db.WithContext(ctx).Raw("SELECT s.team_id, s.points AS points, s.created_at AS timestamp FROM solves s WHERE s.team_id = ? UNION ALL SELECT uh.team_id, -uh.cost AS points, uh.created_at AS timestamp FROM unlocked_hints uh WHERE uh.team_id = ? ORDER BY timestamp;", teamID, teamID).Scan(&rows).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve solve records", err, nil)
	}

	for _, row := range rows {
		result = append(result, row.ToDomain(ctx))
	}

	return result, nil
}

func New(db *gorm.DB) scoreDom.Repository {
	return &scoreRepo{
		db: db,
	}
}
