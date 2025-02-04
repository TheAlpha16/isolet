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

func FetchHint(c *fiber.Ctx, hid int) (models.Hint, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15 * time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := doesHintExist(db, hid); err != nil {
		log.Println(err)
		return models.Hint{}, err
	}

	var existingHint models.Hint
	if err := db.Where("hid = ?", hid).First(&existingHint).Error; err != nil {
		log.Println(err)
		return models.Hint{}, errors.New("failed to fetch challenge")
	}

	return existingHint, nil
}

func SaveHintData(c *fiber.Ctx, hintData *models.Hint) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15 * time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Save(hintData).Where("hid = ?", hintData.HID).Error; err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func doesHintExist(db *gorm.DB, hid int) error {
	if err := db.Where("hid = ?", hid).First(&models.Hint{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("this hint does not exist")
		}
		log.Println(err)
		return errors.New("database error")
	}
	return nil
}
