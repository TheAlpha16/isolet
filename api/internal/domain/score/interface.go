package score

import "context"

type Usecase interface {
	// returns the scoreboard
	GetScoreboard(ctx context.Context, input *GetScoreboardInput) (*Scoreboard, error)

	// returns the scores and solve logs of top 10 teams for graph
	GetScoreGraph(ctx context.Context) (*ScoreGraph, error)

	// returns the count of correct and incorrect submissions for challenges by a team,
	// if there are no submissions for a challenge, submission stats won't exist in the map
	GetSubmissionStats(ctx context.Context, teamID int64, challengeIDs []int64) (map[int64]*SubmissionStats, error)

	// returns the count of solves for challenges
	GetChallengeSolveCounts(ctx context.Context, challengeIDs []int64) (map[int64]int, error)

	// returns the solves of challenges by a team
	GetTeamSolves(ctx context.Context, teamID int64) (map[int64]struct{}, error)

	// returns the score of a team
	GetTeamScore(ctx context.Context, teamID int64) (int, error)

	// applies delta to the score of a team
	// if the call to cache fails, a background job is scheduled to rebuild the scoreboard in cache
	UpdateTeamScore(ctx context.Context, teamID int64, delta int) error

	// refreshes the scoreboard by fetching the latest data from the repository
	RefreshScoreboard(ctx context.Context) error
}

type Repository interface {
	GetScoreboard(ctx context.Context) (map[int64]int, error)
	GetChallengeSolveCounts(ctx context.Context, challengeIDs []int64) (map[int64]int, error)

	GetTeamSolves(ctx context.Context, teamID int64) (map[int64]struct{}, error)
	GetTeamScore(ctx context.Context, teamID int64) (int, error)
	GetSubmissionStats(ctx context.Context, teamID int64, challengeIDs []int64) (map[int64]*SubmissionStats, error)
	GetTeamsSolveRecords(ctx context.Context, teamIDs []int64) (map[int64][]*ScoreRecord, error)
}
