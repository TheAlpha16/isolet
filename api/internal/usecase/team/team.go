package team

import (
	"context"

	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
)

type teamImpl struct {
	repo teamDom.Repository
}

func (t *teamImpl) Create(ctx context.Context, input *teamDom.CreateInput) error {
	return nil
}

func New(repo teamDom.Repository) teamDom.Usecase {
	return &teamImpl{
		repo: repo,
	}
}
