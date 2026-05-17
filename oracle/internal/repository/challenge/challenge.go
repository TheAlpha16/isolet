package challenge

import (
	"context"
	"errors"

	"github.com/TheAlpha16/isolet/oracle/infra/cache"
	challengeDom "github.com/TheAlpha16/isolet/oracle/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	repoCache "github.com/TheAlpha16/isolet/oracle/internal/repository/cache"
	"github.com/TheAlpha16/isolet/oracle/utils"
	"github.com/jackc/pgx/v5/pgconn"

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

			if err := challengeRepo.db.WithContext(ctx).Preload("Category").Preload("Hints").Find(&challenges).Error; err != nil {
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
			if err := challengeRepo.db.WithContext(ctx).Preload("Category").Preload("Hints").First(&challenge, id).Error; err != nil {
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
			Points:       submission.Points,
		}

		if err := tx.Create(solve).Error; err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				if pgErr.Code == errorDom.PgErrDuplicateKey {
					return errorDom.Raise(ctx, errorDom.ErrChallengeAlreadySolved, "", err, nil)
				}
			}
			return errorDom.Raise(ctx, errorDom.ErrDBCreateError, "failed to create solve", err, common.ExtraData{"solve": solve})
		}

		return nil
	})
}

func (challengeRepo *challengeRepo) GetUnlockedHints(ctx context.Context, teamID int64) (map[int64]struct{}, error) {
	result := make(map[int64]struct{}, 100)
	var rows []*UnlockedHint

	if err := challengeRepo.db.WithContext(ctx).Select("hint_id").Where("team_id = ?", teamID).Find(&rows).Error; err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve unlocked hints", err, common.ExtraData{utils.ContextKeyTeamID: teamID})
	}

	for _, row := range rows {
		result[row.HintID] = struct{}{}
	}

	return result, nil
}

func (challengeRepo *challengeRepo) GetHintByID(ctx context.Context, id int64) (*challengeDom.Hint, error) {
	config := utils.GetConfig()

	return repoCache.CachedQuery(
		ctx, challengeRepo.cache,
		GetHintCacheKey(id),
		config.Hints.CacheTTL,
		func() (*challengeDom.Hint, error) {
			var hint Hint
			if err := challengeRepo.db.WithContext(ctx).First(&hint, id).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errorDom.Raise(ctx, errorDom.ErrHintNotFound, "", err, nil)
				}
				return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve hint", err, common.ExtraData{"id": id})
			}

			return hint.ToDomain(ctx)
		},
	)
}

func (challengeRepo *challengeRepo) UnlockHint(ctx context.Context, uHint *challengeDom.UnlockedHint) error {
	uHintModel, err := NewUnlockedHintModel(uHint)
	if err != nil {
		return err
	}

	res := challengeRepo.db.WithContext(ctx).Exec("INSERT INTO unlocked_hints (team_id, hint_id, cost, created_at) SELECT ?, ?, ?, EXTRACT(EPOCH FROM NOW())::bigint WHERE ((SELECT COALESCE(SUM(points), 0) FROM solves WHERE team_id = ?) - (SELECT COALESCE(SUM(cost), 0) FROM unlocked_hints WHERE team_id = ?)) >= ?;", uHintModel.TeamID, uHintModel.HintID, uHintModel.Cost, uHintModel.TeamID, uHintModel.TeamID, uHintModel.Cost)

	if res.Error != nil {
		if pgErr, ok := res.Error.(*pgconn.PgError); ok {
			if pgErr.Code == errorDom.PgErrDuplicateKey {
				return errorDom.Raise(ctx, errorDom.ErrHintAlreadyUnlocked, "", nil, nil)
			}
		}
		return errorDom.Raise(ctx, errorDom.ErrDBExecError, "failed to unlock hint", res.Error, common.ExtraData{utils.ContextKeyTeamID: uHintModel.TeamID, "hint_id": uHintModel.HintID})
	}

	if res.RowsAffected == 0 {
		return errorDom.Raise(ctx, errorDom.ErrHintCostExceeded, "", nil, nil)
	}

	return nil
}

func (challengeRepo *challengeRepo) GetTeamSubmissions(ctx context.Context, teamID int64) ([]*challengeDom.Submission, error) {
	var result []*challengeDom.Submission
	var submissions []*Submission

	err := challengeRepo.db.WithContext(ctx).Where("team_id = ?", teamID).Find(&submissions).Error
	if err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrDBReadError, "failed to retrieve submissions", err, common.ExtraData{utils.ContextKeyTeamID: teamID})
	}

	for _, submission := range submissions {
		domSubmission, err := submission.ToDomain(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, domSubmission)
	}

	return result, nil
}

func New(db *gorm.DB, c cache.Cache) challengeDom.Repository {
	return &challengeRepo{
		db:    db,
		cache: c,
	}
}
