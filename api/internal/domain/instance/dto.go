package instance

import (
	"context"

	"github.com/TheAlpha16/isolet/api/utils/validator"
)

type StartInput struct {
	ChallengeID int64 `json:"challenge_id" validate:"required"`
}

func (si *StartInput) Validate(ctx context.Context) error {
	return validator.Validate(ctx, si)
}
