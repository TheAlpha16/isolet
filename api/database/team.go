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

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT join_team(?, ?, ?)", team.TeamID, user.UserID, config.TEAM_LEN).Error; err != nil {
			return errors.New(CleanSQLException(err.Error()))
		}
		return nil
	})

	if err != nil {
		return errors.New(CleanSQLException(err.Error()))
	}

	if checkUser, err := ReadUser(c, user.UserID); err == nil {
		if checkUser.TeamID != team.TeamID {
			return errors.New("error in joining team")
		}
	} else {
		log.Println(err)
		return errors.New("error in joining team")
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

func GetTeamRank(c *fiber.Ctx, teamid int64) (int64, int64, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var totalTeams int64

	var result struct {
		Rank  int64
		Score int64
	}

	if err := db.Model(&models.Team{}).Count(&totalTeams).Error; err != nil {
		log.Println(err)
		return result.Rank, result.Score, err
	}

	if err := db.Model(&models.Team{}).Select("rank, score").Raw("SELECT rank, score FROM get_scoreboard(?, ?) WHERE teamid = ?", totalTeams, 0, teamid).Scan(&result).Error; err != nil {
		log.Println(err)
		return result.Rank, result.Score, err
	}

	return result.Rank, result.Score, nil
}

func VerifyInviteToken(c *fiber.Ctx, token string) (*models.Team, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var team models.Team
	var tokenData models.Token

	if err := db.Where("token = ?", token).First(&tokenData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid token")
		}
		log.Println(err)
		return nil, errors.New("error in token verification")
	}

	if tokenData.Type != "invite_token" {
		return nil, errors.New("invalid token")
	}

	if tokenData.Expiry.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	if err := db.Joins("JOIN users ON users.userid = ?", tokenData.UserID).Where("teams.teamid = users.teamid").First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid token")
		}
		log.Println(err)
		return nil, errors.New("error in token verification")
	}

	return &team, nil
}

func GenerateInviteToken(c *fiber.Ctx, userid int64) (string, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var token models.Token
	token.UserID = userid
	token.Type = "invite_token"
	token.Expiry = time.Now().Add(30 * time.Minute)

	if err := db.Create(&token).Error; err != nil {
		log.Println(err)
		return "", errors.New("error in token generation")
	}

	return token.Token.String(), nil
}
