package score

import (
	"context"
	"math"

	"github.com/TheAlpha16/isolet/oracle/infra/cache"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	scoreDom "github.com/TheAlpha16/isolet/oracle/internal/domain/score"
	teamDom "github.com/TheAlpha16/isolet/oracle/internal/domain/team"
	"github.com/TheAlpha16/isolet/oracle/utils/logger"

	"go.uber.org/zap"
)

type scoreImpl struct {
	repo   scoreDom.Repository
	teamUc teamDom.Usecase
	cache  cache.Cache
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
	score, err := scoreImpl.cache.ZScore(ctx, scoreDom.ScoreboardCacheKey, scoreDom.ScoreboardMember(teamID))
	if err != nil {
		if !errorDom.Is(err, errorDom.ErrCacheZSetMissingMember) {
			return 0, err
		}
		return 0, nil
	}
	return int(score), nil
}

func (scoreImpl *scoreImpl) UpdateTeamScore(ctx context.Context, teamID int64, delta int) error {
	if err := scoreImpl.cache.ZIncrBy(ctx, scoreDom.ScoreboardCacheKey, float64(delta), scoreDom.ScoreboardMember(teamID)); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrScoreUpdateFailed, "", err, nil)
	}
	return nil
}

func (scoreImpl *scoreImpl) GetScoreboard(ctx context.Context, input *scoreDom.GetScoreboardInput) (*scoreDom.Scoreboard, error) {
	entryLen, err := scoreImpl.cache.ZCard(ctx, scoreDom.ScoreboardCacheKey)
	if err != nil {
		return nil, err
	}

	scoreboard := scoreDom.Scoreboard{
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: int(math.Ceil(float64(entryLen) / float64(input.PageSize))),
		Entries:    []*scoreDom.ScoreboardEntry{},
	}

	if input.Page > scoreboard.TotalPages {
		return &scoreboard, nil
	}

	start := int64((scoreboard.Page - 1) * scoreboard.PageSize)
	stop := start + int64(scoreboard.PageSize) - 1

	items, err := scoreImpl.cache.ZRevRangeWithScores(ctx, scoreDom.ScoreboardCacheKey, start, stop)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return &scoreboard, nil
	}

	entries, _, err := scoreImpl.getScoreboardEntries(ctx, start, stop)
	if err != nil {
		return nil, err
	}
	scoreboard.Entries = entries

	return &scoreboard, nil
}

func (scoreImpl *scoreImpl) GetScoreGraph(ctx context.Context) (*scoreDom.ScoreGraph, error) {
	scoreGraph := scoreDom.ScoreGraph{
		Entries: []*scoreDom.ScoreGraphEntry{},
	}

	var start int64 = 0
	stop := start + int64(scoreDom.ScoreGraphSize) - 1

	entries, teamIDs, err := scoreImpl.getScoreboardEntries(ctx, start, stop)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return &scoreGraph, nil
	}

	scoreRecords, err := scoreImpl.repo.GetTeamsScoreRecords(ctx, teamIDs)
	if err != nil {
		return nil, err
	}

	rows := make([]*scoreDom.ScoreGraphEntry, len(entries))
	for i, entry := range entries {
		graphEntry := &scoreDom.ScoreGraphEntry{
			ScoreboardEntry: *entry,
			Records:         scoreRecords[entry.TeamID],
		}
		if graphEntry.Records == nil {
			graphEntry.Records = []*scoreDom.ScoreRecord{}
		}
		rows[i] = graphEntry
	}
	scoreGraph.Entries = rows
	scoreGraph.Count = len(rows)

	return &scoreGraph, nil
}

func (scoreImpl *scoreImpl) RefreshScoreboard(ctx context.Context) error {
	scores, err := scoreImpl.repo.GetScoreboard(ctx)
	if err != nil {
		return err
	}

	// exit early if there are no scores
	if len(scores) == 0 {
		return nil
	}

	serializedScores := make(map[string]float64, len(scores))
	for teamID, score := range scores {
		serializedScores[scoreDom.ScoreboardMember(teamID)] = float64(score)
	}
	if err := scoreImpl.cache.ZAdd(ctx, scoreDom.ScoreboardCacheKey, serializedScores); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrScoreUpdateFailed, "", err, nil)
	}
	return nil
}

func (scoreImpl *scoreImpl) GetTeamScoreRecords(ctx context.Context, teamID int64) ([]*scoreDom.ScoreRecord, error) {
	return scoreImpl.repo.GetTeamScoreRecords(ctx, teamID)
}

func (scoreImpl *scoreImpl) getScoreboardEntries(ctx context.Context, start, stop int64) ([]*scoreDom.ScoreboardEntry, []int64, error) {
	entries := []*scoreDom.ScoreboardEntry{}

	items, err := scoreImpl.cache.ZRevRangeWithScores(ctx, scoreDom.ScoreboardCacheKey, start, stop)
	if err != nil {
		return nil, nil, err
	}

	if len(items) == 0 {
		return entries, nil, nil
	}

	entries = make([]*scoreDom.ScoreboardEntry, len(items))
	teamIDs := make([]int64, len(items))

	for i, item := range items {
		teamID, err := scoreDom.ParseScoreboardMember(ctx, item.Member)
		if err != nil {
			return nil, nil, err
		}

		entries[i] = &scoreDom.ScoreboardEntry{
			TeamID: teamID,
			Score:  int(item.Score),
			Rank:   int(start) + i + 1,
		}
		teamIDs[i] = teamID
	}

	teamNames, err := scoreImpl.teamUc.GetNameByIDs(ctx, teamIDs)
	if err != nil {
		return nil, nil, err
	}

	for _, entry := range entries {
		entry.TeamName = teamNames[entry.TeamID]
	}

	return entries, teamIDs, nil
}

func New(repo scoreDom.Repository, teamUc teamDom.Usecase, c cache.Cache) scoreDom.Usecase {
	scoreImpl := &scoreImpl{
		repo:   repo,
		teamUc: teamUc,
		cache:  c,
	}
	err := scoreImpl.RefreshScoreboard(context.Background())
	if err != nil {
		errorDom.RaiseToSentry(context.Background(), err)
		logger.GetAppLogger().Fatal("failed to refresh score board", zap.Error(err))
	}
	return scoreImpl
}
