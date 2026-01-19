package team

import (
	"context"
	"strings"

	"github.com/TheAlpha16/isolet/oracle/utils/validator"
)

type CreateInput struct {
	TeamName string `json:"team_name" validate:"required,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func (c *CreateInput) Validate(ctx context.Context) error {
	// normalize
	c.TeamName = strings.TrimSpace(c.TeamName)
	c.Password = strings.TrimSpace(c.Password)

	return validator.Validate(ctx, c)
}

type JoinInput struct {
	TeamName string `json:"team_name" validate:"required,max=32"`
	Password string `json:"password" validate:"required,max=32"`
}

func (j *JoinInput) Validate(ctx context.Context) error {
	// normalize
	j.TeamName = strings.TrimSpace(j.TeamName)
	j.Password = strings.TrimSpace(j.Password)

	return validator.Validate(ctx, j)
}

type GenerateInviteOutput struct {
	InviteLink string `json:"invite_link"`
}
