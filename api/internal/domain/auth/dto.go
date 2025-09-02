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

type RegisterInput struct {
	Username string `json:"username" validate:"required,min=3,max=33"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func (r *RegisterInput) Validate(ctx context.Context) error {
	return validator.Validate(ctx, r)
}

type RegisterOutput struct {
	Session *Session
}
