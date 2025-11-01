package manifest

import (
	"context"

	"github.com/TheAlpha16/isolet/api/infra/cache"
	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
	"gorm.io/gorm"
)

type manifestRepo struct {
	db    *gorm.DB
	cache cache.Cache
}

func (manifestRepo *manifestRepo) GetByChallengeID(ctx context.Context, challengeID int64) (*manifestDom.Manifest, error) {
	return nil, nil
}

func New(db *gorm.DB, cache cache.Cache) manifestDom.Repository {
	return &manifestRepo{
		db:    db,
		cache: cache,
	}
}
