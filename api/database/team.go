package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

func TeamExists(c *fiber.Ctx, teamid int64) (string, bool) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var team models.Team
	if err := db.Where("teamid = ?", teamid).First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", false
		}
		log.Println(err)
		return "", false
	}
	return team.TeamName, true
}

func TeamNameExists(teamname string) bool {
	var team models.Team
	if err := DB.Where("teamname = ?", teamname).First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		log.Println(err)
		return false
	}
	return true
}

func CreateTeam(c *fiber.Ctx, team *models.Team) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Omit("Score", "Rank", "Submissions", "LastSubmission", "Members").Create(team).Error; err != nil {
		return err
	}

	if err := UpdateUserTeam(c, team.Captain, team.TeamID); err != nil {
		return err
	}

	return nil
}

func AuthenticateTeam(c *fiber.Ctx, team *models.Team) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Where("teamname = ? AND password = ?", team.TeamName, team.Password).First(team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("invalid credentials")
		}
		log.Println(err)
		return err
	}

	return nil
}

func JoinTeam(c *fiber.Ctx, user *models.User, team *models.Team) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Raw("SELECT join_team(?, ?, ?)", team.TeamID, user.UserID, config.TEAM_LEN).Error; err != nil {
		return errors.New(CleanSQLException(err.Error()))
	}

	return nil
}

func ReadTeam(c *fiber.Ctx, teamid int64) (models.Team, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var team models.Team
	if err := db.Preload("Members", "teamid = ?",
		teamid,
	).Select("teamid, teamname, captain").Where("teamid = ?", teamid).First(&team).Error; err != nil {
		log.Println(err)
		return team, err
	}

	return team, nil
}

func GetSubmissions(c *fiber.Ctx, teamid int64) ([]models.Sublog, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var submissions []models.Sublog
	if err := db.Select("sublogs.sid, sublogs.chall_id, sublogs.userid, sublogs.teamid, sublogs.correct, sublogs.timestamp, challenges.points AS points").Joins("JOIN challenges ON challenges.chall_id = sublogs.chall_id").Where("teamid = ?", teamid).Find(&submissions).Error; err != nil {
		log.Println(err)
		return submissions, err
	}

	return submissions, nil
}
