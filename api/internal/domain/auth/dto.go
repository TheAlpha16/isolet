package auth

import (
	"context"
	"strings"

	"github.com/TheAlpha16/isolet/api/utils/validator"
)

type LoginInput struct {
	Identifier string `json:"identifier" validate:"required,max=320"`
	Password   string `json:"password" validate:"required,max=32"`
}

func (l *LoginInput) Validate(ctx context.Context) error {
	// normalize
	l.Identifier = strings.TrimSpace(l.Identifier)
	l.Password = strings.TrimSpace(l.Password)

	return validator.Validate(ctx, l)
}

type RegisterInput struct {
	Username string `json:"username" validate:"required,max=32"`
	Email    string `json:"email" validate:"required,email,max=320"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func (r *RegisterInput) Validate(ctx context.Context) error {
	// normalize
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)

	return validator.Validate(ctx, r)
}

type RegisterOutput struct {
	Session *Session
}

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email,max=320"`
}

func (f *ForgotPasswordInput) Validate(ctx context.Context) error {
	// normalize
	f.Email = strings.TrimSpace(strings.ToLower(f.Email))

	return validator.Validate(ctx, f)
}

type ResetPasswordInput struct {
	Token    string `json:"token" validate:"required,max=500"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func (r *ResetPasswordInput) Validate(ctx context.Context) error {
	// normalize
	r.Token = strings.TrimSpace(r.Token)
	r.Password = strings.TrimSpace(r.Password)

	return validator.Validate(ctx, r)
}
