package database

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/api/models"
	"github.com/TheAlpha16/isolet/api/config"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

func CanStartInstance(c *fiber.Ctx, chall_id int, teamid int64) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)
	running := new(models.Running)

	if err := db.Select("runid").Where("chall_id = ? AND teamid = ?", chall_id, teamid).First(running).Error; err == nil {
		return errors.New("instance already running")
	} else if err != gorm.ErrRecordNotFound {
		log.Println(err)
		return errors.New("error in starting the instance, contact admin")
	}

	if err := db.Model(&models.Running{}).Create(&models.Running{ChallID: chall_id, TeamID: teamid}).Error; err != nil {
		if strings.Contains(err.Error(), "start more instances for the team") {
			return errors.New("concurrent instance limit reached for the team")
		}
		log.Println(err)
		return errors.New("error in starting the instance, contact admin")
	}

	return nil
}

func DeleteRunning(c *fiber.Ctx, chall_id int, teamid int64) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Where("chall_id = ? AND teamid = ?", chall_id, teamid).Delete(&models.Running{}).Error; err != nil {
		log.Println(err)
		return errors.New("error in stopping the instance, contact admin")
	}

	return nil
}

func NewFlag(c *fiber.Ctx, flagObject *models.Flag) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Create(flagObject).Error; err != nil {
		log.Println(err)
		return errors.New("error in starting the instance, contact admin")
	}

	return nil
}

func DeleteFlag(c *fiber.Ctx, chall_id int, teamid int64) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)
	
	if err := db.Where("chall_id = ? AND teamid = ?", chall_id, teamid).Delete(&models.Flag{}).Error; err != nil {
		log.Println(err)
		return errors.New("error in stopping the instance, contact admin")
	}

	return nil
}

func ValidOnDemandChallenge(c *fiber.Ctx, chall_id int, teamid int64, challenge *models.Challenge, image *models.Image) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Raw("WITH solved_challenges AS (SELECT ARRAY_AGG(solves.chall_id) AS solved_array FROM solves WHERE teamid = ?) SELECT challenges.type, challenges.flag FROM challenges CROSS JOIN solved_challenges WHERE challenges.chall_id = ? AND challenges.visible = true AND (challenges.requirements = '{}' OR challenges.requirements <@ solved_array)", teamid, chall_id).First(&challenge).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("challenge does not exist")
		}
		log.Println(err)
		return errors.New("error in fetching challenge data")
	}
	
	if challenge.Type != "on-demand" {
		return errors.New("challenge is not on-demand")
	}

	if err := db.Where("chall_id = ?", chall_id).First(image).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Println(err)
		}
		log.Printf("image details not set for %d", chall_id)
		return errors.New("error in starting the instance, contact admin")
	}

	return nil
}

func IsRunning (ctx context.Context, chall_id int, teamid int64) (string, error) {
	var flag string
	db := DB.WithContext(ctx)

	if err := db.Model(&models.Flag{}).
		Select("flag").
		Where("chall_id = ? AND teamid = ?", chall_id, teamid).
    	First(&flag).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Println(err)
		}
		return "", errors.New("instance not running")
	}

	return flag, nil
}

func GetInstances(c *fiber.Ctx, teamid int64) ([]models.Instance, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	instances := make([]models.Instance, 0)

	db := DB.WithContext(ctx)

	if err := db.Model(&models.Flag{}).
		Select("flags.chall_id, flags.password, flags.port, flags.hostname, flags.deadline, images.deployment").
		Joins("JOIN images ON flags.chall_id = images.chall_id").
		Where("flags.teamid = ?", teamid).
		Find(&instances).Error; err != nil {
		log.Println(err)
		return instances, errors.New("error in fetching instances, please contact admin")
	}

	return instances, nil
}

func AddTime(c *fiber.Ctx, chall_id int, teamid int64) (int64, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)
	record := new(models.Flag)

	if err := db.Select("extended, deadline").
		Where("chall_id = ? AND teamid = ?", chall_id, teamid).
		First(&record).Error; err != nil {
		log.Println(err)
		return 0, errors.New("error in extension, please contact admin")
	}

	if (record.Extended + 1) > (config.MAX_INSTANCE_TIME / config.INSTANCE_TIME) {
		return 0, errors.New("limit reached")
	}

	if err := db.Model(&models.Flag{}).
		Where("chall_id = ? AND teamid = ?", chall_id, teamid).
		Update("extended", record.Extended + 1).Error; err != nil {
		log.Println(err)
		return 0, errors.New("error in extension, please contact admin")
	}

	newdeadline := time.UnixMilli(record.Deadline).Add(time.Minute * time.Duration(config.INSTANCE_TIME)).UnixMilli()

	if err := db.Model(&models.Flag{}).
		Where("chall_id = ? AND teamid = ?", chall_id, teamid).
		Update("deadline", newdeadline).Error; err != nil {
		log.Println(err)
		return 0, errors.New("error in extension, please contact admin")
	}

	return newdeadline, nil
}
