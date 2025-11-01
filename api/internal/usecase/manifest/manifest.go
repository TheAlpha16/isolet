package manifest

import (
	"context"

	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
)

type manifestImpl struct {
	repo manifestDom.Repository
}

func (m *manifestImpl) GetByChallengeID(ctx context.Context, challengeID int64) (*manifestDom.Manifest, error) {
	return m.repo.GetByChallengeID(ctx, challengeID)
}

func New(repo manifestDom.Repository) manifestDom.Usecase {
	return &manifestImpl{
		repo: repo,
	}
}
