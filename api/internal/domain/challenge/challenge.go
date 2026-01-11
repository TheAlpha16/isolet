package challenge

import (
	"context"
	"fmt"
	"regexp"

	"github.com/TheAlpha16/isolet/api/internal/domain"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/utils"
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
	Unlocked    bool // this field does not exist in the database
	domain.BaseEntity
}

type UnlockedHint struct {
	ID     int64
	TeamID int64
	HintID int64
	Cost   int
	domain.ImmutableEntity
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

func (c *Challenge) Slug(ctx context.Context) (string, error) {
	slug := utils.SlugifyForSubdomain(c.Name)
	if slug == "" {
		return "", errorDom.Raise(ctx, errorDom.ErrChallengeNameInvalid, "challenge name results in an invalid slug", nil, common.ExtraData{"challenge_name": c.Name, "slug": slug, "challenge_id": c.ID})
	}
	return slug, nil
}

func (c *Challenge) RandomizedFlag(flag string) string {
	re := regexp.MustCompile(utils.RegexGenericFlag)

	matches := re.FindStringSubmatch(flag)
	if len(matches) != 3 {
		return flag
	}

	prefix := matches[1]
	content := matches[2]
	suffix := utils.GenerateRandom()[:utils.InstanceFlagSuffixLength]
	return fmt.Sprintf("%s%s_%s}", prefix, content, suffix)
}

func (h *Hint) IsUnlocked(unlockedHints map[int64]struct{}) bool {
	if _, ok := unlockedHints[h.ID]; ok {
		return true
	}
	return false
}
