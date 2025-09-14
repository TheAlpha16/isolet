package team

import (
	"context"

	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
)

type Usecase interface {
	Create(ctx context.Context, input *CreateInput) (*authDom.Session, error)
}

type Repository interface {
	Create(ctx context.Context, team *Team) (*Team, error)
}
