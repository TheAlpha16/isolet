package challenge

import (
	"context"
	"errors"

	"github.com/TheAlpha16/isolet/api/infra/cache"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	repoCache "github.com/TheAlpha16/isolet/api/internal/repository/cache"
	"github.com/TheAlpha16/isolet/api/utils"

	"gorm.io/gorm"
)

type challengeRepo struct {
	db    *gorm.DB
	cache cache.Cache
}

func (challengeRepo *challengeRepo) GetAll(ctx context.Context) ([]*challengeDom.Challenge, error) {
	config := utils.GetConfig()

	return repoCache.CachedQuery(
		ctx, challengeRepo.cache,
		GetAllChallengesCacheKey(),
		config.Challenges.CacheTTL,
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

func (challengeRepo *challengeRepo) GetByID(ctx context.Context, id int64) (*challengeDom.Challenge, error) {
	config := utils.GetConfig()

	return repoCache.CachedQuery(
		ctx, challengeRepo.cache,
		GetChallengeCacheKey(id),
		config.Challenges.CacheTTL,
		func() (*challengeDom.Challenge, error) {
			var challenge Challenge
			if err := challengeRepo.db.Preload("Category").Preload("Hints").First(&challenge, id).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errorDom.Raise(ctx, errorDom.ErrChallengeNotFound, "", err, nil)
				}
				return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve challenge", err, common.ExtraData{"id": id})
			}

			return challenge.ToDomain(ctx)
		},
	)
}

func (challengeRepo *challengeRepo) SubmitFlag(ctx context.Context, submission *challengeDom.Submission) error {
	submissionModel, err := NewSubmissionModel(submission)
	if err != nil {
		return err
	}

	return challengeRepo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(submissionModel).Error; err != nil {
			return errorDom.Raise(ctx, errorDom.ErrDBCreateError, "failed to create submission", err, common.ExtraData{"submission": submission})
		}

		if !submissionModel.IsCorrect {
			return nil
		}

		solve := &Solve{
			ChallengeID:  submissionModel.ChallengeID,
			TeamID:       submissionModel.TeamID,
			SubmissionID: submissionModel.ID,
		}

		if err := tx.Create(solve).Error; err != nil {
			return errorDom.Raise(ctx, errorDom.ErrDBCreateError, "failed to create solve", err, common.ExtraData{"solve": solve})
		}

		return nil
	})
}

func (challengeRepo *challengeRepo) GetSubmissionStats(ctx context.Context, teamID int64, challengeIDs []int64) (map[int64]*challengeDom.SubmissionStats, error) {
	result := make(map[int64]*challengeDom.SubmissionStats, len(challengeIDs))
	var rows []SubmissionStatsDTO

	if len(challengeIDs) == 0 {
		return result, nil
	}

	if err := challengeRepo.db.WithContext(ctx).Raw("SELECT challenge_id, COUNT(*) FILTER (WHERE is_correct = true)  AS correct_count, COUNT(*) FILTER (WHERE is_correct = false) AS incorrect_count FROM submissions WHERE team_id = ? AND challenge_id IN ? GROUP BY challenge_id", teamID, challengeIDs).Scan(&rows).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve submission stats", err, common.ExtraData{"team_id": teamID, "challenge_ids": challengeIDs})
	}

	for _, row := range rows {
		result[row.ChallengeID] = &challengeDom.SubmissionStats{
			CorrectCount:   int(row.CorrectCount),
			IncorrectCount: int(row.IncorrectCount),
		}
	}

	return result, nil
}

func (challengeRepo *challengeRepo) GetChallengeSolveCounts(ctx context.Context, challengeIDs []int64) (map[int64]int, error) {
	result := make(map[int64]int, len(challengeIDs))
	var rows []ChallengeSolveCountDTO

	if len(challengeIDs) == 0 {
		return result, nil
	}

	if err := challengeRepo.db.WithContext(ctx).Raw("SELECT challenge_id, COUNT(*) AS solve_count FROM solves WHERE challenge_id IN ? GROUP BY challenge_id", challengeIDs).Scan(&rows).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve challenge solve counts", err, common.ExtraData{"challenge_ids": challengeIDs})
	}

	for _, row := range rows {
		result[row.ChallengeID] = int(row.SolveCount)
	}

	return result, nil
}

func New(db *gorm.DB, cache cache.Cache) challengeDom.Repository {
	return &challengeRepo{
		db:    db,
		cache: cache,
	}
}
