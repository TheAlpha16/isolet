package challenge

import (
	"github.com/TheAlpha16/isolet/api/infra/database/postgres"
	"github.com/lib/pq"
)

type Category struct {
	postgres.BaseModel
	Name      string `gorm:"uniqueIndex;not null"`
	IsVisible bool   `gorm:"not null;default:false"`
}

type Hint struct {
	postgres.BaseModel
	Text        string `gorm:"not null"`
	Cost        int    `gorm:"not null;default:0"`
	IsVisible   bool   `gorm:"not null;default:false"`
	ChallengeID int64  `gorm:"column:challenge_id;not null;index"`

	Challenge Challenge `gorm:"foreignKey:ChallengeID;references:ID"`
}

type Challenge struct {
	postgres.BaseModel
	Name        string         `gorm:"uniqueIndex;not null"`
	Prompt      string         `gorm:"not null"`
	CategoryID  int64          `gorm:"column:category_id;not null"`
	Flag        string         `gorm:"type:text"`
	Type        string         `gorm:"type:challenge_type;not null;default:'static'"`
	Points      int            `gorm:"not null"`
	Files       pq.StringArray `gorm:"type:text[];default:'{}'"`
	Author      string         `gorm:"not null;default:anonymous"`
	Tags        pq.StringArray `gorm:"type:text[];default:'{}'"`
	Links       pq.StringArray `gorm:"type:text[];default:'{}'"`
	IsVisible   bool           `gorm:"not null;default:false"`
	MaxAttempts int            `gorm:"not null;default:0"`

	Category Category `gorm:"foreignKey:CategoryID;references:ID"`
	Hints []Hint `gorm:"foreignKey:ChallengeID;references:ID;constraint:OnDelete:CASCADE"`
}
