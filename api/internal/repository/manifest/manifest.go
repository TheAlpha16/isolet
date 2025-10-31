package manifest

import (
	"github.com/TheAlpha16/isolet/api/infra/database/postgres"
	challengeRepo "github.com/TheAlpha16/isolet/api/internal/repository/challenge"
)

type Manifest struct {
	postgres.BaseModel
	ChallengeID int64      `gorm:"not null"`
	Image       string     `gorm:"not null"`
	Flag        *string    `gorm:"type:text"`
	Requests    []Resource `gorm:"foreignKey:ManifestID;references:ID;constraint:OnDelete:CASCADE"`
	Limits      []Resource `gorm:"foreignKey:ManifestID;references:ID;constraint:OnDelete:CASCADE"`
	Endpoints   []Endpoint `gorm:"foreignKey:ManifestID;references:ID;constraint:OnDelete:CASCADE"`

	Challenge challengeRepo.Challenge `gorm:"foreignKey:ChallengeID;references:ID;constraint:OnDelete:CASCADE"`
}

type Resource struct {
	postgres.BaseModel
	Name       string `gorm:"type:resource_name_type;not null"`
	Value      string `gorm:"not null"`
	Type       string `gorm:"type:resource_type;not null"`
	ManifestID int64  `gorm:"not null;index"`

	Manifest Manifest `gorm:"foreignKey:ManifestID;references:ID"`
}

type Endpoint struct {
	postgres.BaseModel
	Name       string `gorm:"not null"`
	Protocol   string `gorm:"type:protocol_type;not null"`
	TargetPort int32  `gorm:"not null"`
	ManifestID int64  `gorm:"not null;index"`

	Manifest Manifest `gorm:"foreignKey:ManifestID;references:ID"`
}
