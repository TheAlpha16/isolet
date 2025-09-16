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
	ID   int64  `json:"id"`
	Text string `json:"text"`
	Cost int    `json:"cost"`
}

type ChallengeDTO struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Prompt      string        `json:"prompt"`
	Category    CategoryDTO   `json:"category"`
	Type        ChallengeType `json:"type"`
	Points      int           `json:"points"`
	Files       []string      `json:"files"`
	Hints       []*HintDTO    `json:"hints"`
	Author      string        `json:"author"`
	Tags        []string      `json:"tags"`
	Links       []string      `json:"links"`
	MaxAttempts int           `json:"max_attempts"`
}

func (cat *Category) ToDTO() *CategoryDTO {
	return &CategoryDTO{
		ID:   cat.ID,
		Name: cat.Name,
	}
}

func (h *Hint) ToDTO() *HintDTO {
	return &HintDTO{
		ID:   h.ID,
		Text: h.Text,
		Cost: h.Cost,
	}
}

func (ch *Challenge) ToDTO() *ChallengeDTO {
	var hints []*HintDTO
	for _, hint := range ch.Hints {
		hints = append(hints, hint.ToDTO())
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

type SubmissionStats struct {
	CorrectCount   int `json:"correct_count"`
	IncorrectCount int `json:"incorrect_count"`
}
