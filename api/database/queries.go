package database

import (
	"log"
	"time"
	"errors"
	"context"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddToChallenges(chall models.Challenge) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	// Use GORM's Upsert equivalent for "ON CONFLICT DO UPDATE"
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "level"}},
		DoUpdates: clause.AssignmentColumns([]string{"chall_name", "prompt", "tags"}),
	}).Create(&chall).Error

	if err != nil {
		return err
	}

	return nil
}

func ReadChallenges(c *fiber.Ctx, teamid int64) (map[string][]models.Challenge, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var team models.Team
	if err := db.Where("teamid = ?", teamid).First(&team).Error; err != nil {
		return nil, err
	}

	var challenges []models.Challenge
	if err := db.Preload("Category").
		Preload("Hints", "visible = ?", true).
		Where("visible = ?", true).
		Find(&challenges).Error; err != nil {
		return nil, err
	}

	// Post-fetch filtering and modifications
	filteredChallenges := make(map[string][]models.Challenge)

	for _, challenge := range challenges {
		requirementsMet := isRequirementMet(challenge.Requirements, team.Solved)
		if !requirementsMet {
			continue
		}

		for i, hint := range challenge.Hints {
			hintUnlocked := isHintUnlocked(int64(hint.HID), team.UHints)
			if hint.Cost > 0 && !hintUnlocked {
				challenge.Hints[i].Hint = ""
			}

			challenge.Hints[i].Unlocked = hintUnlocked
		}

		if challenge.Type == "dynamic" {
			imageMetaData := new(models.Image) 
			
			if err := db.Select("deployment, port, subd").Where("chall_id = ?", challenge.ChallID).First(imageMetaData).Error; err != nil {
				return nil, err
			}
			
			connLink := GenerateChallengeEndpoint(imageMetaData.Deployment, imageMetaData.Subd, config.INSTANCE_HOSTNAME, imageMetaData.Port)

			challenge.Links = append(challenge.Links, connLink)
		}

		challenge.Done = isChallengeSolved(int64(challenge.ChallID), team.Solved)

		if catChallenges, exists := filteredChallenges[challenge.Category.CategoryName]; exists {
			filteredChallenges[challenge.Category.CategoryName] = append(catChallenges, challenge)
		} else {
			filteredChallenges[challenge.Category.CategoryName] = []models.Challenge{challenge}
		}
	}

	return filteredChallenges, nil
}

func ValidFlagEntry(ctx context.Context, chall_id int, teamid int64) (models.Challenge, error) {
	db := DB.WithContext(ctx)
	var err error

	var challenge models.Challenge
	if err := db.Select("chall_name, type, flag, points, requirements").Where("chall_id = ? AND visible = ?", chall_id, true).First(&challenge).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Println(err)
		}
		return challenge, errors.New("challenge does not exist")
	}

	var team models.Team
	if err := db.Select("solved").Where("teamid = ?", teamid).First(&team).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return challenge, errors.New("team does not exist")
		}

		log.Println(err)
		return challenge, errors.New("error in fetching team data")
	}

	if !isRequirementMet(challenge.Requirements, team.Solved) {
		return challenge, errors.New("challenge requirements not met")
	}

	if isChallengeSolved(int64(chall_id), team.Solved) {
		return challenge, errors.New("challenge already solved")
	}

	if challenge.Type == "on-demand" {
		if challenge.Flag, err = IsRunning(ctx, chall_id, teamid); err != nil {
			return challenge, err
		}
	}

	return challenge, nil
}

func VerifyFlag(c *fiber.Ctx, chall_id int, userid int64, teamid int64, flag string) (bool, string) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	challenge, err := ValidFlagEntry(ctx, chall_id, teamid)
	if err != nil {
		return false, err.Error()
	}

	db := DB.WithContext(ctx)

	sublog := models.Sublog{
		ChallID: chall_id,
		UserID:  userid,
		TeamID:  teamid,
		Flag:    flag,
		Correct: flag == challenge.Flag,
		IP: c.Locals("clientIP").(string),
	}

	if config.POST_EVENT != "false" {
		if !sublog.Correct {
			return sublog.Correct, "incorrect flag"
		} else {
			return sublog.Correct, "correct flag"
		}
	}

	if err := db.Create(&sublog).Error; err != nil {
		log.Println(err)
		return false, "error in verification, please contact admin"
	}

	if !sublog.Correct {
		return false, "incorrect flag"
	}

	if err := db.Model(&models.Team{}).
		Where("teamid = ?", teamid).
		Update("solved", gorm.Expr("array_append(solved, ?)", chall_id)).
		Error; err != nil {
		log.Println(err)
		return false, "error in verification, please contact admin"
	}

	if err := db.Model(&models.User{}).
		Where("userid = ?", userid).
		Update("score", gorm.Expr("score + ?", challenge.Points)).
		Error; err != nil {
		log.Println(err)
		return false, "error in verification, please contact admin"
	}

	if err := db.Model(&models.Challenge{}).
		Where("chall_id = ?", chall_id).
		Update("solves", gorm.Expr("solves + 1")).
		Error; err != nil {
		log.Println(err)
		return false, "error in verification, please contact admin"
	}

	return true, "correct flag"
}

func ReadScores(c *fiber.Ctx, page int) (models.ScoreBoard, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	perPage := 10 
	offset := (page - 1) * perPage
	scores := make([]models.Score, 0)
	var totalTeams int64
	
	if err := db.Model(&models.Team{}).Count(&totalTeams).Error; err != nil {
		return models.ScoreBoard{}, err
	}

	pageCount := int((totalTeams + int64(perPage) - 1) / int64(perPage))
	if pageCount < page {
		return models.ScoreBoard{
			PageCount: pageCount,
			Page:      page,
			Scores:    scores,
		}, nil
	}

	err := db.Table("users").
		Select("teams.teamid as teamid, teams.teamname as teamname, SUM(users.score) as score").
		Joins("LEFT JOIN teams ON users.teamid = teams.teamid").
		Group("teams.teamid").
		Order("SUM(users.score) DESC").
		Limit(perPage).
		Offset(offset).
		Scan(&scores).Error

	if err != nil {
		return models.ScoreBoard{}, err
	}

	for i := range scores {
		scores[i].Rank = offset + i + 1
	}

	return models.ScoreBoard{
		PageCount: pageCount,
		Page:      page,
		Scores:    scores,
	}, nil
}
