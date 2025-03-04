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

func FetchConfig(c *fiber.Ctx, key string) (models.Config, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)
	
	if err := doesConfigExist(db, key); err != nil {
		log.Println(err)
		return models.Config{}, err
	}

	var existingConfig models.Config
	if err := db.Where("key = ?", key).First(&existingConfig).Error; err != nil {
		log.Println(err)
		return models.Config{}, errors.New("failed to fetch challenge")
	}

	return existingConfig, nil
}


func SaveConfigData(c *fiber.Ctx, updatedConfig *models.Config) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Save(updatedConfig).Where("key = ?", updatedConfig.Key).Error; err != nil {
		log.Println(err)
		return errors.New("failed to save updated challenge")
	}

	return nil
}


func doesConfigExist(db *gorm.DB, key string) error {
	if err := db.Where("key = ?", key).First(&models.Config{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("this config does not exist")
		}
		log.Println(err)
		return errors.New("database error")
	}

	return nil
}