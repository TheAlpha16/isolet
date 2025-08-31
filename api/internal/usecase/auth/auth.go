package auth

import (
	"context"
	"time"

	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
)

type AuthUsecase struct {
	repo authDom.Repository
}

func (a *AuthUsecase) Login(ctx context.Context, input *authDom.LoginInput) (*authDom.Session, error) {
	return &authDom.Session{
		ID:        "<session_id>",
		UserID:    6969,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

func New(repo authDom.Repository) authDom.Usecase {
	return &AuthUsecase{
		repo: repo,
	}
}
