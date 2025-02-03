package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/TheAlpha16/isolet/admin/models"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

func FetchChallenge(c *fiber.Ctx, challid int) (models.Challenge, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)
	
	if err := doesChallengeExist(db, challid); err != nil {
		log.Println(err)
		return models.Challenge{}, err
	}

	var existingChallenge models.Challenge
	if err := db.Where("chall_id = ?", challid).First(&existingChallenge).Error; err != nil {
		log.Println(err)
		return models.Challenge{}, errors.New("failed to fetch challenge")
	}

	return existingChallenge, nil
}


func SaveChallengeMetaData(c *fiber.Ctx, updatedChall *models.Challenge) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Save(updatedChall).Where("chall_id = ?", updatedChall.ChallID).Error; err != nil {
		log.Println(err)
		return errors.New("failed to save updated challenge")
	}

	return nil
}

func doesChallengeExist(db *gorm.DB, challID int) error {
	if err := db.Where("chall_id = ?", challID).First(&models.Challenge{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("this challenge does not exist")
		}
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}

