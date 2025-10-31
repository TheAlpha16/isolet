package instance

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/api/infra/database/postgres"
	"github.com/TheAlpha16/isolet/api/internal/domain"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
	challengeRepo "github.com/TheAlpha16/isolet/api/internal/repository/challenge"
	teamRepo "github.com/TheAlpha16/isolet/api/internal/repository/team"
)

type Instance struct {
	postgres.BaseModel
	TeamID         *int64  `gorm:"column:team_id"`
	ChallengeID    int64   `gorm:"not null"`
	ExpiresAt      *int64  `gorm:"column:expires_at"`
	AvailableAt    *int64  `gorm:"column:available_at"`
	AllowExtension bool    `gorm:"column:allow_extension;not null;default:true"`
	Flag           *string `gorm:"type:text"`

	Challenge challengeRepo.Challenge `gorm:"foreignKey:ChallengeID;references:ID;constraint:OnDelete:CASCADE"`
	Team      teamRepo.Team           `gorm:"foreignKey:TeamID;references:ID;constraint:OnDelete:CASCADE"`
}

func (in *Instance) ToDomain(ctx context.Context) (*instanceDom.Instance, error) {
	lifecycle := instanceDom.Lifecycle{}
	if in.ExpiresAt != nil {
		expiresAt := time.Unix(*in.ExpiresAt, 0)
		lifecycle.ExpiresAt = &expiresAt
	}
	if in.AvailableAt != nil {
		availableAt := time.Unix(*in.AvailableAt, 0)
		lifecycle.AvailableAt = &availableAt
	}
	lifecycle.AllowExtension = in.AllowExtension

	return &instanceDom.Instance{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(in.CreatedAt, 0),
			UpdatedAt: time.Unix(in.UpdatedAt, 0),
		},
		ID: in.ID,
		Manifest: &manifestDom.Manifest{
			ChallengeID: in.ChallengeID,
			Flag:        in.Flag,
		},
		TeamID:    in.TeamID,
		Lifecycle: &lifecycle,
	}, nil
}

func NewInstanceModel(in *instanceDom.Instance) (*Instance, error) {
	instance := &Instance{
		ChallengeID: in.Manifest.ChallengeID,
		TeamID:      in.TeamID,
		Flag:        in.Manifest.Flag,
	}
	if in.ID != 0 {
		instance.ID = in.ID
	}

	if in.Lifecycle != nil {
		if in.Lifecycle.ExpiresAt != nil {
			expiresAt := in.Lifecycle.ExpiresAt.Unix()
			instance.ExpiresAt = &expiresAt
		}
		if in.Lifecycle.AvailableAt != nil {
			availableAt := in.Lifecycle.AvailableAt.Unix()
			instance.AvailableAt = &availableAt
		}
		instance.AllowExtension = in.Lifecycle.AllowExtension
	}

	return instance, nil
}
