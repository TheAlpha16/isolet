package challenge

import (
	"context"

	"github.com/TheAlpha16/isolet/api/infra/cache"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	repoCache "github.com/TheAlpha16/isolet/api/internal/repository/cache"
	"github.com/TheAlpha16/isolet/api/utils"

	"gorm.io/gorm"
)

const (
	cacheKeyChallenges = "challenges"
)

type challengeRepo struct {
	db    *gorm.DB
	cache cache.Cache
}

func (challengeRepo *challengeRepo) GetAll(ctx context.Context) ([]*challengeDom.Challenge, error) {
	config := utils.GetConfig()

	return repoCache.CachedQuery(
		ctx, challengeRepo.cache, cacheKeyChallenges, config.Challenges.CacheTTL,
		func() ([]*challengeDom.Challenge, error) {
			var challenges []*Challenge

			if err := challengeRepo.db.Preload("Category").Preload("Hints").Find(&challenges).Error; err != nil {
				return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve challenges", err, nil)
			}

			var domChallenges []*challengeDom.Challenge
			for _, challenge := range challenges {
				domChallenge, err := challenge.ToDomain(ctx)
				if err != nil {
					return nil, err
				}
				domChallenges = append(domChallenges, domChallenge)
			}

			return domChallenges, nil
		},
	)
}

func New(db *gorm.DB, cache cache.Cache) challengeDom.Repository {
	return &challengeRepo{
		db:    db,
		cache: cache,
	}
}
