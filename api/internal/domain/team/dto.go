package team

import (
	"context"

	"github.com/TheAlpha16/isolet/api/utils/validator"
)

type CreateInput struct {
	TeamName string `json:"team_name" validate:"required,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func (c *CreateInput) Validate(ctx context.Context) error {
	return validator.Validate(ctx, c)
}

type JoinInput struct {
	TeamName string `json:"team_name" validate:"required,max=32"`
	Password string `json:"password" validate:"required,max=32"`
}

func (j *JoinInput) Validate(ctx context.Context) error {
	return validator.Validate(ctx, j)
}
