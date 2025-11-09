package manifest

import (
	"time"

	"github.com/TheAlpha16/isolet/api/infra/database/postgres"
	"github.com/TheAlpha16/isolet/api/internal/domain"
	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
	challengeRepo "github.com/TheAlpha16/isolet/api/internal/repository/challenge"
)

type Manifest struct {
	postgres.BaseModel
	ChallengeID   int64          `gorm:"not null;uniqueIndex"`
	Image         string         `gorm:"not null"`
	FlagTemplate  *string        `gorm:"type:text"`
	Requests      []Resource     `gorm:"foreignKey:ManifestID;references:ID;constraint:OnDelete:CASCADE"`
	Limits        []Resource     `gorm:"foreignKey:ManifestID;references:ID;constraint:OnDelete:CASCADE"`
	EndpointSpecs []EndpointSpec `gorm:"foreignKey:ManifestID;references:ID;constraint:OnDelete:CASCADE"`

	Challenge challengeRepo.Challenge `gorm:"foreignKey:ChallengeID;references:ID;constraint:OnDelete:CASCADE"`
}

func (m *Manifest) ToDomain() *manifestDom.Manifest {
	requests := make([]*manifestDom.Resource, len(m.Requests))
	for i, req := range m.Requests {
		requests[i] = req.ToDomain()
	}

	limits := make([]*manifestDom.Resource, len(m.Limits))
	for i, lim := range m.Limits {
		limits[i] = lim.ToDomain()
	}

	endpointSpecs := make([]*manifestDom.EndpointSpec, len(m.EndpointSpecs))
	for i, ep := range m.EndpointSpecs {
		endpointSpecs[i] = ep.ToDomain()
	}

	return &manifestDom.Manifest{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(m.CreatedAt, 0),
			UpdatedAt: time.Unix(m.UpdatedAt, 0),
		},
		ID:            m.ID,
		ChallengeID:   m.ChallengeID,
		Image:         m.Image,
		FlagTemplate:  m.FlagTemplate,
		Requests:      requests,
		Limits:        limits,
		EndpointSpecs: endpointSpecs,
	}
}

type Resource struct {
	postgres.BaseModel
	Name       string `gorm:"type:resource_name_type;not null"`
	Value      string `gorm:"not null"`
	Type       string `gorm:"type:resource_type;not null"`
	ManifestID int64  `gorm:"not null;index"`

	Manifest Manifest `gorm:"foreignKey:ManifestID;references:ID"`
}

func (res *Resource) ToDomain() *manifestDom.Resource {
	return &manifestDom.Resource{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(res.CreatedAt, 0),
			UpdatedAt: time.Unix(res.UpdatedAt, 0),
		},
		ID:         res.ID,
		Name:       manifestDom.ResourceName(res.Name),
		Value:      res.Value,
		Type:       manifestDom.ResourceType(res.Type),
		ManifestID: res.ManifestID,
	}
}

type EndpointSpec struct {
	postgres.BaseModel
	Name       string `gorm:"not null"`
	Protocol   string `gorm:"type:protocol_type;not null"`
	TargetPort int32  `gorm:"not null"`
	ManifestID int64  `gorm:"not null;index"`

	Manifest Manifest `gorm:"foreignKey:ManifestID;references:ID"`
}

func (ep *EndpointSpec) ToDomain() *manifestDom.EndpointSpec {
	return &manifestDom.EndpointSpec{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(ep.CreatedAt, 0),
			UpdatedAt: time.Unix(ep.UpdatedAt, 0),
		},
		ID:         ep.ID,
		Name:       ep.Name,
		Protocol:   manifestDom.Protocol(ep.Protocol),
		TargetPort: ep.TargetPort,
		ManifestID: ep.ManifestID,
	}
}
