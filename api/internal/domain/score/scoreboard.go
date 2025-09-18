package score

import (
	"context"
	"strconv"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
)

const (
	maxPageSize     = 100
	defaultPageSize = 50

	ScoreboardCacheKey      = "scoreboard"
	TeamScoreCacheKeyPrefix = "score:team"
)

type SubmissionStats struct {
	CorrectCount   int `json:"correct_count"`
	IncorrectCount int `json:"incorrect_count"`
}

func ScoreboardMember(teamID int64) string {
	return strconv.FormatInt(teamID, 10)
}

func ParseScoreboardMember(ctx context.Context, member string) (int64, error) {
	teamID, err := strconv.ParseInt(member, 10, 64)
	if err != nil {
		return 0, errorDom.RaiseInternal(ctx, "invalid team id in scoreboard member", err, common.ExtraData{"member": member})
	}
	return teamID, nil
}

type Scoreboard struct {
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
	Entries    []*ScoreboardEntry `json:"entries"`
}

type ScoreboardEntry struct {
	TeamID   int64  `json:"team_id"`
	TeamName string `json:"team_name"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}
