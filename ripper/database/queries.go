package database

import (
	"context"
	"log"
	"time"

	"github.com/TheAlpha16/isolet/ripper/models"
)

func DeleteFlag(teamid int64, chall_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Where("teamid = ? AND chall_id = ?", teamid, chall_id).Delete(&models.Flag{}).Error; err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DeleteRunning(teamid int64, chall_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Where("teamid = ? AND chall_id = ?", teamid, chall_id).Delete(&models.Running{}).Error; err != nil {
		log.Println(err)
		return err
	}

	return nil
}