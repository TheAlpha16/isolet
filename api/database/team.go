package database

import (
	"context"
	"errors"
	"time"
	"log"

	"github.com/TheAlpha16/isolet/api/models"
	"github.com/TheAlpha16/isolet/api/config"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

func TeamExists(teamname string) bool {
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

	if err := db.Create(team).Error; err != nil {
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

func JoinTeam(c *fiber.Ctx, teamname string, userid int) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)
	team := new(models.Team)

	if err := db.Where("teamname = ?", teamname).First(team).Error; err != nil {
		log.Println(err)
		return errors.New("unable to join the team. contact admin")
	}

	if len(team.Members) >= config.TEAM_LEN {
		return errors.New("team is full")
	}

	team.Members = append(team.Members, int64(userid))

	if err := db.Model(&team).Update("members", team.Members).Error; err != nil {
		log.Println(err)
		return errors.New("unable to join the team. contact admin")
	}

	if err := UpdateUserTeam(c, userid, team.TeamID); err != nil {
		log.Println(err)
		return errors.New("unable to join the team. contact admin")
	}

	return nil
}

