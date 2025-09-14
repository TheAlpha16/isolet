package auth

import (
	"context"

	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type Usecase interface {
	Login(ctx context.Context, input *LoginInput) (*Session, error)
	Register(ctx context.Context, input *RegisterInput) (*Session, error)
	Verify(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, input *ForgotPasswordInput) error
	ResetPassword(ctx context.Context, input *ResetPasswordInput) error
	
	GenerateAuthToken(ctx context.Context, user *userDom.User) (*tokenDom.Token, string, error)
}
