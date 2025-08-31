package auth

import (
	"context"

	"github.com/TheAlpha16/isolet/api/utils/validator"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func (l *LoginInput) Validate(ctx context.Context) error {
	return validator.Validate(ctx, l)
}
