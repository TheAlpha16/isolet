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
}
