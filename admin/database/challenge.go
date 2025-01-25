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
func EditChallengeMetaData(c *fiber.Ctx, updatedChall *models.Challenge) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var existingChall models.Challenge
	if err := doesChallengeExist(db, updatedChall.ChallID, existingChall); err != nil {
		log.Println(err)
		return err
	}

	updates := map[string]interface{}{
		"chall_name":    updatedChall.Name,
		"prompt":        updatedChall.Prompt,
		"category_name": updatedChall.Category.CategoryName,
		"type":          updatedChall.Type,
		"points":        updatedChall.Points,
		"flag":          updatedChall.Flag,
		"author":        updatedChall.Author,
		"tags":          updatedChall.Tags,
		"visible":       updatedChall.Visible,
		"links":         updatedChall.Links,
	}

	if err := performChallengeTransaction(db, existingChall, updates); err != nil {
		log.Println(err)
		return err
	}
	
	return nil
}

func doesChallengeExist(db *gorm.DB, challID int, challenge models.Challenge) error {
	if err := db.Where("chall_id = ?", challID).First(&challenge).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("this challenge does not exist")
		}
		return errors.New("database error")
	}

	return nil
}

func performChallengeTransaction(db *gorm.DB, existingChallenge models.Challenge, updates map[string]interface{}) error {

	tx := db.Begin()
	if tx.Error != nil {
		return errors.New("error in starting a transaction")
	}

	if err := tx.Model(&existingChallenge).Updates(updates).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update challenge")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("failed to commit changes")
	}

	return nil
}
