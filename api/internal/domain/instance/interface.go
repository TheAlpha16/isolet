package instance

import "context"

type Usecase interface {
	Start(ctx context.Context, input *StartInput) (*InstanceDTO, error)
	Stop(ctx context.Context, input *StopInput) error
	Extend(ctx context.Context, input *ExtendInput) (*InstanceDTO, error)
}

type Repository interface {
	Create(ctx context.Context, instance *Instance) (*Instance, error)
	GetByID(ctx context.Context, id int64) (*Instance, error)
	Update(ctx context.Context, instance *Instance, fields []string) error
	Delete(ctx context.Context, id int64) error
	GetByTeamAndChallenge(ctx context.Context, teamID, challengeID int64) (*Instance, error)
}

type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Extend(ctx context.Context) error
}
