package instance

import (
	"context"

	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
)

type Usecase interface {
	List(ctx context.Context) ([]*InstanceDTO, error)
	Start(ctx context.Context, input *StartInput) (*InstanceDTO, error)
	Stop(ctx context.Context, input *StopInput) error
	Extend(ctx context.Context, input *ExtendInput) (*InstanceDTO, error)
	HandleEvent(ctx context.Context, eventType string, instance *Instance) error
}

type Repository interface {
	Create(ctx context.Context, instance *Instance) (*Instance, error)
	GetByID(ctx context.Context, id int64) (*Instance, error)
	GetByRefs(ctx context.Context, teamID *int64, challengeID int64) (*Instance, error)
	GetByTeam(ctx context.Context, teamID int64) ([]*Instance, error)
	Update(ctx context.Context, instance *Instance, fields []string) error
	Delete(ctx context.Context, id int64) error
	DeleteByRefs(ctx context.Context, teamID *int64, challengeID int64) error
}

type Service interface {
	Start(ctx context.Context, inst *Instance, manifest *manifestDom.Manifest) error
	Stop(ctx context.Context, inst *Instance) error
	Extend(ctx context.Context, inst *Instance) error
}
