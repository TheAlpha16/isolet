package auth

import "context"

type Repository interface{}

type Usecase interface {
	Login(ctx context.Context, input *LoginInput) (*Session, error)
	Register(ctx context.Context, input *RegisterInput) (*Session, error)
	Verify(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, input *ForgotPasswordInput) error
}
