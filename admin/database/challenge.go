package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/TheAlpha16/isolet/admin/models"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddToChallenges(chall models.Challenge) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	// Use GORM's Upsert equivalent for "ON CONFLICT DO UPDATE"
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "level"}},
		DoUpdates: clause.AssignmentColumns([]string{"chall_name", "prompt", "tags"}),
	}).Create(&chall).Error

	if err != nil {
		return err
	}

	return nil
}

// Admin functions
func EditChallengeData(c *fiber.Ctx, challengeData *models.ChallengeData) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var existingChallData models.ChallengeData
	if err := db.Where("chall_id = ?", challengeData.ChallID).First(&existingChallData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println(err)
			return errors.New("this challenge does not exist")
		}
		log.Println(err)
		return errors.New("database error")
	}

	updates := map[string]interface{}{
		"chall_name":    challengeData.Name,
		"prompt":        challengeData.Prompt,
		"category_name": challengeData.CategoryName,
		"type":          challengeData.Type,
		"points":        challengeData.Points,
		"files":         challengeData.Files,
		"author":        challengeData.Author,
		"port":          challengeData.Port,
		"tags":          challengeData.Tags,
		"links":         challengeData.Links,
	}

	tx := db.Begin()
	if tx.Error != nil {
		log.Println(tx.Error)
		return errors.New("error in starting a transaction")
	}

	if err := tx.Model(&existingChallData).Updates(updates).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return errors.New("failed to update challenge")
	}

	// updating hints separately
	if len(challengeData.Hints) > 0 {
		if err := UpdateHints(tx, challengeData.ChallID, challengeData.Hints); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return errors.New("failed to commit changes")
	}
	return nil
}
