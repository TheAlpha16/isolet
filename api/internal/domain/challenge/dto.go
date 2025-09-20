package challenge

import (
	"context"
	"strings"

	"github.com/TheAlpha16/isolet/api/utils/validator"
)

type CategoryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type HintDTO struct {
	ID          int64  `json:"id"`
	Text        string `json:"text"`
	Cost        int    `json:"cost"`
	Unlocked    bool   `json:"unlocked"`
	ChallengeID int64  `json:"challenge_id"`
}

type ChallengeDTO struct {
	ID           int64         `json:"id"`
	Name         string        `json:"name"`
	Prompt       string        `json:"prompt"`
	Category     CategoryDTO   `json:"category"`
	Type         ChallengeType `json:"type"`
	Points       int           `json:"points"`
	Files        []string      `json:"files"`
	Hints        []*HintDTO    `json:"hints"`
	Author       string        `json:"author"`
	Tags         []string      `json:"tags"`
	Links        []string      `json:"links"`
	MaxAttempts  int           `json:"max_attempts"`
	TotalSolves  int           `json:"total_solves"`
	Solved       bool          `json:"solved"`
	AttemptCount int           `json:"attempt_count"`
}

type SubmissionDTO struct {
	ChallengeID int64 `json:"challenge_id"`
	UserID      int64 `json:"user_id"`
	TeamID      int64 `json:"team_id,omitempty"`
	IsCorrect   bool  `json:"is_correct"`
	Timestamp   int64 `json:"timestamp"`
}

func (cat *Category) ToDTO() *CategoryDTO {
	return &CategoryDTO{
		ID:   cat.ID,
		Name: cat.Name,
	}
}

func (h *Hint) ToDTO(challengeID int64) *HintDTO {
	return &HintDTO{
		ID:          h.ID,
		Text:        h.Text,
		Cost:        h.Cost,
		Unlocked:    h.Unlocked,
		ChallengeID: challengeID,
	}
}

func (ch *Challenge) ToDTO() *ChallengeDTO {
	var hints []*HintDTO
	for _, hint := range ch.Hints {
		if !hint.IsVisible {
			continue
		}
		hints = append(hints, hint.ToDTO(ch.ID))
	}

	if len(hints) == 0 {
		hints = []*HintDTO{}
	}

	return &ChallengeDTO{
		ID:          ch.ID,
		Name:        ch.Name,
		Prompt:      ch.Prompt,
		Category:    *ch.Category.ToDTO(),
		Type:        ch.Type,
		Points:      ch.Points,
		Files:       ch.Files,
		Hints:       hints,
		Author:      ch.Author,
		Tags:        ch.Tags,
		Links:       ch.Links,
		MaxAttempts: ch.MaxAttempts,
	}
}

func (sub *Submission) ToDTO() *SubmissionDTO {
	return &SubmissionDTO{
		ChallengeID: sub.ChallengeID,
		UserID:      sub.UserID,
		TeamID:      sub.TeamID,
		IsCorrect:   sub.IsCorrect,
		Timestamp:   sub.CreatedAt.Unix(),
	}
}

type SubmitFlagInput struct {
	ChallengeID int64  `json:"challenge_id" validate:"required"`
	Flag        string `json:"flag" validate:"required"`
}

func (sfi *SubmitFlagInput) Validate(ctx context.Context) error {
	// normalize
	sfi.Flag = strings.TrimSpace(sfi.Flag)

	return validator.Validate(ctx, sfi)
}

type SubmitFlagOutput struct {
	IsCorrect bool `json:"is_correct"`
}

type UnlockHintInput struct {
	HintID int64 `json:"hint_id" validate:"required"`
}

func (sfi *UnlockHintInput) Validate(ctx context.Context) error {
	return validator.Validate(ctx, sfi)
}
