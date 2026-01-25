package instance

import (
	"context"

	"github.com/TheAlpha16/isolet/oracle/utils/validator"
)

type InstanceDTO struct {
	ID          int64          `json:"id"`
	ChallengeID int64          `json:"challenge_id"`
	TeamID      int64          `json:"team_id,omitempty"`
	ExpiresAt   int64          `json:"expires_at,omitempty"`
	Endpoints   []*EndpointDTO `json:"endpoints"`
}

func (in *Instance) ToDTO() *InstanceDTO {
	inst := &InstanceDTO{
		ID:          in.ID,
		ChallengeID: in.ChallengeID,
	}
	if in.Lifecycle.ExpiresAt != nil {
		inst.ExpiresAt = in.Lifecycle.ExpiresAt.Unix()
	}
	if in.TeamID != nil {
		inst.TeamID = *in.TeamID
	}
	if len(in.Endpoints) > 0 {
		inst.Endpoints = make([]*EndpointDTO, len(in.Endpoints))
		for i, ep := range in.Endpoints {
			inst.Endpoints[i] = ep.ToDTO()
		}
	}
	return inst
}

type EndpointDTO struct {
	Name     string  `json:"name"`
	Protocol string  `json:"protocol"`
	Hostname *string `json:"hostname,omitempty"`
	Port     *int32  `json:"port,omitempty"`
	Ready    bool    `json:"ready"`
}

func (ep *Endpoint) ToDTO() *EndpointDTO {
	return &EndpointDTO{
		Name:     ep.Name,
		Protocol: string(ep.Protocol),
		Hostname: ep.Hostname,
		Port:     ep.Port,
		Ready:    ep.Ready,
	}
}

type StartInput struct {
	ChallengeID int64 `json:"challenge_id" validate:"required"`
}

func (si *StartInput) Validate(ctx context.Context) error {
	return validator.Validate(ctx, si)
}

type StopInput struct {
	InstanceID int64 `json:"instance_id" validate:"required"`
}

func (si *StopInput) Validate(ctx context.Context) error {
	return validator.Validate(ctx, si)
}

type ExtendInput struct {
	InstanceID int64 `json:"instance_id" validate:"required"`
}

func (ei *ExtendInput) Validate(ctx context.Context) error {
	return validator.Validate(ctx, ei)
}
