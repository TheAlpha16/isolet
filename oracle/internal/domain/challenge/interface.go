package challenge

import "context"

type Usecase interface {
	List(ctx context.Context) ([]*ChallengeDTO, error)
	ValidateAccess(ctx context.Context, challengeID int64, teamID int64) (*Challenge, error)
	SubmitFlag(ctx context.Context, input *SubmitFlagInput) (*SubmitFlagOutput, error)
	UnlockHint(ctx context.Context, input *UnlockHintInput) (*HintDTO, error)
	GetTeamSubmissions(ctx context.Context, teamID int64) ([]*Submission, error)
}

type Repository interface {
	// challenges
	GetAll(ctx context.Context) ([]*Challenge, error)
	GetByID(ctx context.Context, id int64) (*Challenge, error)

	// submissions + solves
	SubmitFlag(ctx context.Context, submission *Submission) error
	GetTeamSubmissions(ctx context.Context, teamID int64) ([]*Submission, error)

	// hints
	GetHintByID(ctx context.Context, id int64) (*Hint, error)
	GetUnlockedHints(ctx context.Context, teamID int64) (map[int64]struct{}, error)
	UnlockHint(ctx context.Context, uHint *UnlockedHint) error
}
