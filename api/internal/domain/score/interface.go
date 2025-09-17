package score

import "context"

type Usecase interface {
	// returns the count of correct and incorrect submissions for challenges by a team,
	// if there are no submissions for a challenge, submission stats won't exist in the map
	GetSubmissionStats(ctx context.Context, teamID int64, challengeIDs []int64) (map[int64]*SubmissionStats, error)

	// returns the count of solves for challenges
	GetChallengeSolveCounts(ctx context.Context, challengeIDs []int64) (map[int64]int, error)
}

type Repository interface {
	GetSubmissionStats(ctx context.Context, teamID int64, challengeIDs []int64) (map[int64]*SubmissionStats, error)
	GetChallengeSolveCounts(ctx context.Context, challengeIDs []int64) (map[int64]int, error)
}
