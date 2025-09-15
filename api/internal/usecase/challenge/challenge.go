package challenge

import (
	"context"

	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
)

type challengeImpl struct {
	repo challengeDom.Repository
}

func (c *challengeImpl) List(ctx context.Context) ([]*challengeDom.Challenge, error) {
	return c.repo.GetAll(ctx)
}

func New(repo challengeDom.Repository) challengeDom.Usecase {
	return &challengeImpl{
		repo: repo,
	}
}
