package instance

import (
	"context"

	"github.com/TheAlpha16/isolet/api/utils/validator"
)

type InstanceDTO struct {
	ID          int64 `json:"id"`
	ChallengeID int64 `json:"challenge_id"`
	TeamID      int64 `json:"team_id"`
	ExpiresAt   int64 `json:"expires_at"`
}

func (in *Instance) ToDTO() *InstanceDTO {
	return &InstanceDTO{
		ID:          in.ID,
		ChallengeID: in.ChallengeID,
		TeamID:      in.TeamID,
		ExpiresAt:   in.ExpiresAt,
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
