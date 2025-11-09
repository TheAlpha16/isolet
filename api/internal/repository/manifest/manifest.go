package manifest

import (
	"context"
	"errors"

	"github.com/TheAlpha16/isolet/api/infra/cache"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
	repoCache "github.com/TheAlpha16/isolet/api/internal/repository/cache"
	"github.com/TheAlpha16/isolet/api/utils"

	"gorm.io/gorm"
)

type manifestRepo struct {
	db    *gorm.DB
	cache cache.Cache
}

func (manifestRepo *manifestRepo) GetByChallengeID(ctx context.Context, challengeID int64) (*manifestDom.Manifest, error) {
	config := utils.GetConfig()
	extraData := common.ExtraData{"challenge_id": challengeID}

	return repoCache.CachedQuery(
		ctx, manifestRepo.cache,
		GetManifestCacheKey(challengeID),
		config.Manifest.CacheTTL,
		func() (*manifestDom.Manifest, error) {
			var manifest Manifest
			if err := manifestRepo.db.WithContext(ctx).
				Preload("Requests", "type = ?", manifestDom.ResourceTypeRequest).
				Preload("Limits", "type = ?", manifestDom.ResourceTypeLimit).
				Preload("EndpointSpecs").
				Where("challenge_id = ?", challengeID).
				First(&manifest).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errorDom.Raise(ctx, errorDom.ErrManifestNotFound, "", err, extraData)
				}
				return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve manifest", err, extraData)
			}
			return manifest.ToDomain(), nil
		},
	)
}

func New(db *gorm.DB, cache cache.Cache) manifestDom.Repository {
	return &manifestRepo{
		db:    db,
		cache: cache,
	}
}
