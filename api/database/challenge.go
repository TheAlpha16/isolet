package database

import (
	"log"
	"time"
	"errors"
	"context"
	"strings"

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

func ReadChallenges(c *fiber.Ctx, teamid int64) (map[string][]models.ChallengeData, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var challenges []models.ChallengeData
	if err := db.Raw("SELECT chall_id, chall_name, prompt, type, points, files, hints, solves, author, tags, links, category_name, deployment, port, subd, done FROM get_challenges(?)", teamid).Scan(&challenges).Error; err != nil {
		return nil, err
	}

	// Post-fetch filtering and modifications
	filteredChallenges := make(map[string][]models.ChallengeData)

	for _, challenge := range challenges {
		if challenge.Type == "dynamic" {
			connLink := GenerateChallengeEndpoint(challenge.Deployment, challenge.Subd, config.INSTANCE_HOSTNAME, challenge.Port)

			challenge.Links = append(challenge.Links, connLink)
		}

		if catChallenges, exists := filteredChallenges[challenge.CategoryName]; exists {
			filteredChallenges[challenge.CategoryName] = append(catChallenges, challenge)
		} else {
			filteredChallenges[challenge.CategoryName] = []models.ChallengeData{challenge}
		}
	}

	return filteredChallenges, nil
}

func ValidFlagEntry(ctx context.Context, chall_id int, teamid int64) (models.Challenge, error) {
	db := DB.WithContext(ctx)
	var err error

	var challenge models.Challenge
	if err := db.Raw("WITH solved_challenges AS (SELECT ARRAY_AGG(solves.chall_id) AS solved_array FROM solves WHERE teamid = ?) SELECT challenges.type, challenges.flag, challenges.chall_id = any(solved_array) AS done FROM challenges CROSS JOIN solved_challenges WHERE challenges.chall_id = ? AND challenges.visible = true AND (challenges.requirements = '{}' OR challenges.requirements <@ solved_array)", teamid, chall_id).Scan(&challenge).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return challenge, errors.New("challenge does not exist")
		}
		log.Println(err)
		return challenge, errors.New("error in fetching challenge data")
	}

	if challenge.Done {
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

	return true, "correct flag"
}

func UnlockHint(c *fiber.Ctx, hid int, teamid int64) (bool, string) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var hint string
	var error_msg string
	if err := db.Raw("SELECT unlock_hint(?, ?)", teamid, hid).Scan(&hint).Error; err != nil {
		error_msg = err.Error()
		if strings.Contains(error_msg, "hint ") || strings.Contains(error_msg, "insufficient") {
			error_msg = strings.TrimPrefix(error_msg, "ERROR: ")
			prefixIndex := strings.Index(error_msg, " (SQLSTATE")
			if prefixIndex != -1 {
				error_msg = error_msg[:prefixIndex]
			}

			return false, error_msg
		}

		log.Println(err)
		return false, "error in unlocking hint"
	}

	return true, hint
}
