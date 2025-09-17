package challenge

import "context"

type Usecase interface {
	List(ctx context.Context) ([]*ChallengeDTO, error)
	ValidateAttempt(ctx context.Context, challengeID int64, teamID int64) (*Challenge, error)
	SubmitFlag(ctx context.Context, input *SubmitFlagInput) (*SubmitFlagOutput, error)
}

type Repository interface {
	// returns all the challenges
	GetAll(ctx context.Context) ([]*Challenge, error)

	// returns a challenge by id
	GetByID(ctx context.Context, id int64) (*Challenge, error)

	// adds an entry to the submissions table
	// and an entry to solves if the submission is correct
	SubmitFlag(ctx context.Context, submission *Submission) error

	// returns the count of correct and incorrect submissions for challenges by a team,
	// if there are no submissions for a challenge, submission stats won't exist in the map
	GetSubmissionStats(ctx context.Context, teamID int64, challengeIDs []int64) (map[int64]*SubmissionStats, error)

	// returns the count of solves for challenges
	GetChallengeSolveCounts(ctx context.Context, challengeIDs []int64) (map[int64]int, error)
}
