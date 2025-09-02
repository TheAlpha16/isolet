package auth

import (
	"context"
	"time"

	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
)

type authImpl struct {
	repo authDom.Repository
}

func (a *authImpl) Login(ctx context.Context, input *authDom.LoginInput) (*authDom.Session, error) {
	return &authDom.Session{
		UserID:    6969,
		Token:     "some-token-here",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

func (a *authImpl) Register(ctx context.Context, input *authDom.RegisterInput) error {
	return nil
}

func New(repo authDom.Repository) authDom.Usecase {
	return &authImpl{
		repo: repo,
	}
}
