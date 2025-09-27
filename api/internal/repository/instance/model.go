package instance

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/api/infra/database/postgres"
	"github.com/TheAlpha16/isolet/api/internal/domain"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	challengeRepo "github.com/TheAlpha16/isolet/api/internal/repository/challenge"
	teamRepo "github.com/TheAlpha16/isolet/api/internal/repository/team"
)

type Instance struct {
	postgres.BaseModel
	TeamID      int64 `gorm:"not null;uniqueIndex:idx_instances_team_challenge"`
	ChallengeID int64 `gorm:"not null;uniqueIndex:idx_instances_team_challenge"`
	ExpiresAt   int64 `gorm:"not null"`

	Challenge challengeRepo.Challenge `gorm:"foreignKey:ChallengeID;references:ID;constraint:OnDelete:CASCADE"`
	Team      teamRepo.Team           `gorm:"foreignKey:TeamID;references:ID;constraint:OnDelete:CASCADE"`
}

func (in *Instance) ToDomain(ctx context.Context) (*instanceDom.Instance, error) {
	return &instanceDom.Instance{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(in.CreatedAt, 0),
			UpdatedAt: time.Unix(in.UpdatedAt, 0),
		},
		ID:          in.ID,
		ChallengeID: in.ChallengeID,
		TeamID:      in.TeamID,
		ExpiresAt:   in.ExpiresAt,
	}, nil
}

func NewInstanceModel(in *instanceDom.Instance) (*Instance, error) {
	return &Instance{
		ChallengeID: in.ChallengeID,
		TeamID:      in.TeamID,
		ExpiresAt:   in.ExpiresAt,
	}, nil
}
