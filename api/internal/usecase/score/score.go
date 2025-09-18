package score

import (
	"context"
	"fmt"
	"math"

	"github.com/TheAlpha16/isolet/api/infra/cache"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	scoreDom "github.com/TheAlpha16/isolet/api/internal/domain/score"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
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
	return scoreImpl.repo.GetTeamScore(ctx, teamID)
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

	// DEBUG
	fmt.Println("entryLen:", entryLen)
	fmt.Println("page", input.Page)
	fmt.Println("pageSize", input.PageSize)
	fmt.Println("total_pages", scoreboard.TotalPages)

	if input.Page > scoreboard.TotalPages {
		return &scoreboard, nil
	}

	start := int64((scoreboard.Page - 1) * scoreboard.PageSize)
	stop := start + int64(scoreboard.PageSize) + 1

	// DEBUG
	fmt.Println("start", start)
	fmt.Println("stop", stop)

	items, err := scoreImpl.cache.ZRevRangeWithScores(ctx, scoreDom.ScoreboardCacheKey, start, stop)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return &scoreboard, nil
	}

	// DEBUG
	fmt.Println("items", items)
	fmt.Println()

	entries := make([]*scoreDom.ScoreboardEntry, len(items))
	teamIDs := make([]int64, len(items))

	for i, item := range items {
		teamID, err := scoreDom.ParseScoreboardMember(ctx, item.Member)
		if err != nil {
			return nil, err
		}
		// DEBUG
		fmt.Println("teamID", teamID)
		fmt.Println("score", item.Score)

		entries[i] = &scoreDom.ScoreboardEntry{
			TeamID: teamID,
			Score:  int(item.Score),
			Rank:   i + 1,
		}
		teamIDs[i] = teamID
	}

	// DEBUG
	fmt.Println("teamIDs", teamIDs)

	teamNames, err := scoreImpl.teamUc.GetNameByIDs(ctx, teamIDs)
	if err != nil {
		return nil, err
	}

	// DEBUG
	fmt.Println("teamNames", teamNames)

	for _, entry := range entries {
		entry.TeamName = teamNames[entry.TeamID]
	}
	scoreboard.Entries = entries

	return &scoreboard, nil
}

func (scoreImpl *scoreImpl) rebuildScoreboard(ctx context.Context) error {
	scores, err := scoreImpl.repo.GetScoreboard(ctx)
	if err != nil {
		return err
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

func New(repo scoreDom.Repository, teamUc teamDom.Usecase, cache cache.Cache) scoreDom.Usecase {
	return &scoreImpl{
		repo:   repo,
		teamUc: teamUc,
		cache:  cache,
	}
}
