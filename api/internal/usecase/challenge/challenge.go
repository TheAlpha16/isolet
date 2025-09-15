package challenge

import (
	"context"

	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
)

type challengeImpl struct {
	repo challengeDom.Repository
}

func (c *challengeImpl) List(ctx context.Context) ([]*challengeDom.ChallengeDTO, error) {
	domChallenges, err := c.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var challenges []*challengeDom.ChallengeDTO
	for _, challenge := range domChallenges {
		challenges = append(challenges, challenge.ToDTO())
	}

	return challenges, nil
}

func New(repo challengeDom.Repository) challengeDom.Usecase {
	return &challengeImpl{
		repo: repo,
	}
}
