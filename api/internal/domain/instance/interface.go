package instance

import "context"

type Usecase interface {
	Start(ctx context.Context, input *StartInput) (*Instance, error)
}

type Repository interface {
	Create(ctx context.Context, instance *Instance) (*Instance, error)
	GetByTeamAndChallenge(ctx context.Context, teamID, challengeID int64) (*Instance, error)
}
