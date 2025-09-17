package challenge

import (
	"github.com/TheAlpha16/isolet/api/internal/domain"
)

type ChallengeType string

const (
	ChallengeStatic   ChallengeType = "static"
	ChallengeDynamic  ChallengeType = "dynamic"
	ChallengeOnDemand ChallengeType = "on-demand"
)

type Category struct {
	ID        int64
	Name      string
	IsVisible bool
	domain.BaseEntity
}

type Hint struct {
	ID          int64
	Text        string
	Cost        int
	IsVisible   bool
	ChallengeID int64
	domain.BaseEntity
}

type Challenge struct {
	ID           int64
	Name         string
	Prompt       string
	Category     Category
	Flag         string
	Type         ChallengeType
	Points       int
	Requirements []int64
	Files        []string
	Hints        []*Hint
	Author       string
	Tags         []string
	Links        []string
	IsVisible    bool
	MaxAttempts  int
	domain.BaseEntity
}

func (c *Challenge) AreRequirementsMet(solves map[int64]struct{}) bool {
	for _, requirement := range c.Requirements {
		if _, ok := solves[requirement]; !ok {
			return false
		}
	}
	return true
}
