package instance

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/oracle/infra/database/postgres"
	"github.com/TheAlpha16/isolet/oracle/internal/domain"
	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/oracle/internal/domain/manifest"
	challengeRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/challenge"
	teamRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/team"
)

type Instance struct {
	postgres.BaseModel
	TeamID         *int64     `gorm:"column:team_id"`
	ChallengeID    int64      `gorm:"not null"`
	Flag           *string    `gorm:"type:text"`
	ExpiresAt      *int64     `gorm:"column:expires_at"`
	AvailableAt    *int64     `gorm:"column:available_at"`
	AllowExtension bool       `gorm:"column:allow_extension;not null;default:true"`
	Endpoints      []Endpoint `gorm:"foreignKey:InstanceID;references:ID;constraint:OnDelete:CASCADE"`

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

	endpoints := make([]*instanceDom.Endpoint, len(in.Endpoints))
	for i, ep := range in.Endpoints {
		endpoints[i] = ep.ToDomain()
	}

	return &instanceDom.Instance{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(in.CreatedAt, 0),
			UpdatedAt: time.Unix(in.UpdatedAt, 0),
		},
		ID:          in.ID,
		TeamID:      in.TeamID,
		ChallengeID: in.ChallengeID,
		Flag:        in.Flag,
		Lifecycle:   &lifecycle,
		Endpoints:   endpoints,
	}, nil
}

func NewInstanceModel(in *instanceDom.Instance) (*Instance, error) {
	instance := &Instance{
		ChallengeID: in.ChallengeID,
		TeamID:      in.TeamID,
		Flag:        in.Flag,
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

	if len(in.Endpoints) > 0 {
		instance.Endpoints = make([]Endpoint, len(in.Endpoints))
		for i, ep := range in.Endpoints {
			instance.Endpoints[i] = *NewEndpointModel(ep)
		}
	}

	return instance, nil
}

type Endpoint struct {
	postgres.BaseModel
	InstanceID int64  `gorm:"not null;index"`
	Name       string `gorm:"not null"`
	Protocol   string `gorm:"type:protocol_type;not null"`
	TargetPort int32  `gorm:"not null"`
	Hostname   string `gorm:"type:text"`
	Port       *int32 `gorm:"column:port"`
	Ready      bool   `gorm:"default:false;not null"`

	Instance Instance `gorm:"foreignKey:InstanceID;references:ID"`
}

func (ep *Endpoint) ToDomain() *instanceDom.Endpoint {
	return &instanceDom.Endpoint{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(ep.CreatedAt, 0),
			UpdatedAt: time.Unix(ep.UpdatedAt, 0),
		},
		ID:         ep.ID,
		InstanceID: ep.InstanceID,
		Name:       ep.Name,
		Protocol:   manifestDom.Protocol(ep.Protocol),
		TargetPort: ep.TargetPort,
		Hostname:   ep.Hostname,
		Port:       ep.Port,
		Ready:      ep.Ready,
	}
}

func NewEndpointModel(ep *instanceDom.Endpoint) *Endpoint {
	return &Endpoint{
		InstanceID: ep.InstanceID,
		Name:       ep.Name,
		Protocol:   string(ep.Protocol),
		TargetPort: ep.TargetPort,
		Hostname:   ep.Hostname,
		Port:       ep.Port,
		Ready:      ep.Ready,
	}
}
