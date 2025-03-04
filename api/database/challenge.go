package database

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

func ReadChallenges(c *fiber.Ctx, teamid int64) (map[string][]models.Challenge, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var challenges []models.Challenge
	if err := db.Raw("SELECT chall_id, chall_name, prompt, type, points, files, hints, solves, author, tags, links, category_name, deployment, port, subd, done, attempts, sub_count FROM get_challenges(?)", teamid).Scan(&challenges).Error; err != nil {
		return nil, err
	}

	// Post-fetch filtering and modifications
	filteredChallenges := make(map[string][]models.Challenge)
	instanceHostname, err := config.Get("INSTANCE_HOSTNAME")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, challenge := range challenges {
		if challenge.Type == "dynamic" {
			connLink := GenerateChallengeEndpoint(challenge.Deployment, challenge.Subd, instanceHostname, challenge.Port)

			challenge.Links = append(challenge.Links, connLink)
		}

		if catChallenges, exists := filteredChallenges[challenge.CategoryName]; exists {
			filteredChallenges[challenge.CategoryName] = append(catChallenges, challenge)
		} else {
			filteredChallenges[challenge.CategoryName] = []models.Challenge{challenge}
		}
	}

	return filteredChallenges, nil
}

func ValidFlagEntry(ctx context.Context, chall_id int, teamid int64) (models.Challenge, error) {
	db := DB.WithContext(ctx)
	var err error

	var challenge models.Challenge
	if err := db.Raw("WITH solved_challenges AS (SELECT ARRAY_AGG(solves.chall_id) AS solved_array FROM solves WHERE teamid = ?) SELECT challenges.type, challenges.flag, challenges.chall_id = any(solved_array) AS done, challenges.attempts, COALESCE((SELECT COUNT(*)::integer FROM sublogs WHERE sublogs.chall_id = ? AND sublogs.teamid = ?), 0) AS sub_count FROM challenges CROSS JOIN solved_challenges WHERE challenges.chall_id = ? AND challenges.visible = true AND (challenges.requirements = '{}' OR challenges.requirements <@ solved_array)", teamid, chall_id, teamid, chall_id).First(&challenge).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return challenge, errors.New("challenge does not exist")
		}
		log.Println(err)
		return challenge, errors.New("error in fetching challenge data")
	}

	if challenge.Done {
		return challenge, errors.New("challenge already solved")
	}

	if challenge.SubCount >= challenge.Attempts {
		return challenge, errors.New("maximum attempts reached")
	}

	if challenge.Type == "on-demand" {
		if challenge.Flag, err = IsRunning(ctx, chall_id, teamid); err != nil {
			return challenge, err
		}
	}

	return challenge, nil
}

func VerifyFlag(c *fiber.Ctx, chall_id int, userid int64, teamid int64, flag string) (bool, string, int) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	challenge, err := ValidFlagEntry(ctx, chall_id, teamid)
	if err != nil {
		return false, err.Error(), -1
	}

	db := DB.WithContext(ctx)

	sublog := models.Sublog{
		ChallID: chall_id,
		UserID:  userid,
		TeamID:  teamid,
		Flag:    flag,
		Correct: flag == challenge.Flag,
		IP:      c.Locals("clientIP").(string),
	}

	postEvent, err := config.GetBool("POST_EVENT")
	if err != nil {
		log.Println(err)
		return false, "error in verification, please contact admin", -1
	}

	if postEvent {
		if !sublog.Correct {
			return sublog.Correct, "incorrect flag", -1
		} else {
			return sublog.Correct, "correct flag", -1
		}
	}

	if sublog.Correct {
		if err := db.Omit("Points", "Timestamp").Create(&sublog).Error; err != nil {
			log.Println(err)
			return false, "error in verification, please contact admin", -1
		}
	}

	if !sublog.Correct {
		return false, "incorrect flag", challenge.SubCount + 1
	}

	return true, "correct flag", challenge.SubCount + 1
}

func UnlockHint(c *fiber.Ctx, hid int, teamid int64) (bool, string) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var hint string
	if err := db.Raw("SELECT unlock_hint(?, ?)", teamid, hid).Scan(&hint).Error; err != nil {
		if strings.Contains(err.Error(), "hint ") || strings.Contains(err.Error(), "insufficient") {
			return false, CleanSQLException(err.Error())
		}

		log.Println(err)
		return false, "error in unlocking hint"
	}

	return true, hint
}

func GetChallengeName(c *fiber.Ctx, chall_id int) (string, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var name string
	if err := db.Model(&models.Challenge{}).Select("chall_name").Where("chall_id = ?", chall_id).First(&name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return name, errors.New("challenge does not exist")
		}
		log.Println(err)
		return name, errors.New("error in fetching challenge name")
	}

	return name, nil
}
