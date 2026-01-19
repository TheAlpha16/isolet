package manifest

import (
	"context"
)

type Usecase interface {
	GetByChallengeID(ctx context.Context, challengeID int64) (*Manifest, error)
}

type Repository interface {
	GetByChallengeID(ctx context.Context, challengeID int64) (*Manifest, error)
}
