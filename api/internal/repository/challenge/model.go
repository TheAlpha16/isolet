package challenge

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/api/infra/database/postgres"
	"github.com/TheAlpha16/isolet/api/internal/domain"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
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
	Hints    []Hint   `gorm:"foreignKey:ChallengeID;references:ID;constraint:OnDelete:CASCADE"`
}

func (cat *Category) ToDomain(ctx context.Context) (*challengeDom.Category, error) {
	return &challengeDom.Category{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(cat.CreatedAt, 0),
			UpdatedAt: time.Unix(cat.UpdatedAt, 0),
		},
		ID:        cat.ID,
		Name:      cat.Name,
		IsVisible: cat.IsVisible,
	}, nil
}

func (h *Hint) ToDomain(ctx context.Context) (*challengeDom.Hint, error) {
	return &challengeDom.Hint{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(h.CreatedAt, 0),
			UpdatedAt: time.Unix(h.UpdatedAt, 0),
		},
		ID:          h.ID,
		Text:        h.Text,
		Cost:        h.Cost,
		IsVisible:   h.IsVisible,
		ChallengeID: h.ChallengeID,
	}, nil
}

func (c *Challenge) ToDomain(ctx context.Context) (*challengeDom.Challenge, error) {
	category, err := c.Category.ToDomain(ctx)
	if err != nil {
		return nil, err
	}

	var hints []*challengeDom.Hint
	for _, hint := range c.Hints {
		hintDomain, err := hint.ToDomain(ctx)
		if err != nil {
			return nil, err
		}
		hints = append(hints, hintDomain)
	}

	return &challengeDom.Challenge{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(c.CreatedAt, 0),
			UpdatedAt: time.Unix(c.UpdatedAt, 0),
		},
		ID:          c.ID,
		Name:        c.Name,
		Prompt:      c.Prompt,
		Category:    *category,
		Flag:        c.Flag,
		Type:        challengeDom.ChallengeType(c.Type),
		Points:      c.Points,
		Files:       c.Files,
		Hints:       hints,
		Author:      c.Author,
		Tags:        c.Tags,
		Links:       c.Links,
		IsVisible:   c.IsVisible,
		MaxAttempts: c.MaxAttempts,
	}, nil
}
