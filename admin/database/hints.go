package database

import (
	"errors"
	"log"

	"github.com/TheAlpha16/isolet/admin/models"
	"gorm.io/gorm"
)

func UpdateHints(tx *gorm.DB, challID int, newHints []models.Hint) error {

	// Get existing hints
	var existingHints []models.Hint
	if err := tx.Where("chall_id = ?", challID).Find(&existingHints).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return errors.New("failed to fetch existing hints")
	}

	existingHintMap := make(map[int]models.Hint)
	for _, hint := range existingHints {
		existingHintMap[hint.HID] = hint
	}

	// Checking if any of the hints are already unlocked
	var unlockedHints []models.UHint
	if err := tx.Where("hid IN ?", getHIDs(existingHints)).Find(&unlockedHints).Error; err != nil {
		log.Println(err)
		return errors.New("failed to get unlocked hints")
	}

	unlockedHintMap := make(map[int]bool)
	for _, uHint := range unlockedHints {
		unlockedHintMap[uHint.HID] = true
	}

	// Update hints
	for _, newHint := range newHints {
		newHint.ChallID = challID

		if existingHint, exists := existingHintMap[newHint.HID]; exists {
			
			// preserve unlocked hints
			if unlockedHintMap[newHint.HID] {
				continue
			}

			// update locked hints
			updates := map[string]interface{} {
				"hint": newHint.Hint,
				"cost": newHint.Cost,
				"visible": newHint.Visible,
			}

			if err := tx.Model(&existingHint).Updates(updates).Error; err != nil {
				log.Println(err)
				return errors.New("failed to update hint")
			}
		} else {
			if err := tx.Create(&newHint).Error; err != nil {
				log.Println(err)
				return errors.New("failed to create new hint")
			}
		}
	}

	return nil
}

func getHIDs(hints []models.Hint) []int {
	hids := make([]int, len(hints))
    for i, hint := range hints {
        hids[i] = hint.HID
    }
    return hids
}