package database

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/api/models"

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
	var team models.Team

	if err := db.Select("type, flag, requirements").Where("chall_id = ? AND visible = ?", chall_id, true).First(challenge).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Println(err)
		}
		return errors.New("challenge does not exist")
	}
	
	if challenge.Type != "on-demand" {
		return errors.New("challenge is not on-demand")
	}

	if err := db.Select("solved").Where("teamid = ?", teamid).First(&team).Error; err != nil {
		log.Println(err)
		return errors.New("error in fetching team data")
	}

	if !isRequirementMet(challenge.Requirements, team.Solved) {
		return errors.New("challenge does not exist")
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

// func GetInstances(c *fiber.Ctx, userid int) ([]models.Instance, error) {
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	instances := make([]models.Instance, 0)
// 	rows, err := DB.QueryContext(ctx, `SELECT userid, level, password, port, verified, hostname, deadline from flags WHERE userid = $1`, userid)
// 	if err != nil {
// 		return instances, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		instance := new(models.Instance)
// 		if err := rows.Scan(&instance.UserID, &instance.Level, &instance.Password, &instance.Port, &instance.Verified, &instance.Hostname, &instance.Deadline); err != nil {
// 			return instances, err
// 		}
// 		instances = append(instances, *instance)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return instances, err
// 	}
// 	return instances, nil
// }


// func AddTime(c *fiber.Ctx, userid int, level int) (bool, string, int64) {
// 	var current int
// 	var deadline int64
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if err := DB.QueryRowContext(ctx, `SELECT extended, deadline FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&current, &deadline); err != nil {
// 		log.Println(err)
// 		return false, "error in extension, please contact admin", 1
// 	}

// 	if (current + 1) > (config.MAX_INSTANCE_TIME / config.INSTANCE_TIME) {
// 		return false, "limit reached", 1
// 	}

// 	_, err := DB.QueryContext(ctx, `UPDATE flags SET extended = $1 WHERE userid = $2 AND level = $3`, current+1, userid, level)
// 	if err != nil {
// 		log.Println(err)
// 		return false, "error in extension, please contact admin", 1
// 	}

// 	newdeadline := time.UnixMilli(deadline).Add(time.Minute * time.Duration(config.INSTANCE_TIME)).UnixMilli()

// 	_, err = DB.QueryContext(ctx, `UPDATE flags SET deadline = $1 WHERE userid = $2 AND level = $3`, newdeadline, userid, level)
// 	if err != nil {
// 		log.Println(err)
// 		return false, "error in extension, please contact admin", 1
// 	}
// 	return true, "", newdeadline
// }
