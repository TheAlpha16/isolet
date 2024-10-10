package database

import (
	"context"
	"time"

	// "github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"

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

func ReadChallenges(c *fiber.Ctx, teamid int) (map[string][]models.Challenge, error) {
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
			
			// fetch deployment, port, subd from images table using chall_id as key
			if err := db.Select("deployment, port, subd").Where("chall_id = ?", challenge.ChallID).First(imageMetaData).Error; err != nil {
				return nil, err
			}
			
			connLink := GenerateChallengeEndpoint(imageMetaData.Deployment, imageMetaData.Subd, imageMetaData.Port)

			challenge.Links = append(challenge.Links, connLink)
		}

		challenge.Done = isChallengeSolved(challenge.Name, team.Solved)

		if catChallenges, exists := filteredChallenges[challenge.Category.CategoryName]; exists {
			filteredChallenges[challenge.Category.CategoryName] = append(catChallenges, challenge)
		} else {
			filteredChallenges[challenge.Category.CategoryName] = []models.Challenge{challenge}
		}
	}

	return filteredChallenges, nil
}

// func CanStartInstance(c *fiber.Ctx, userid int, level int) bool {
// 	var runid int
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if err := DB.QueryRowContext(ctx, `SELECT runid FROM running WHERE userid = $1 AND level = $2`, userid, level).Scan(&runid); err == nil {
// 		return false
// 	}

// 	if _, err := DB.QueryContext(ctx, `INSERT INTO running (userid, level) VALUES ($1, $2)`, userid, level); err != nil {
// 		log.Println(err)
// 		return false
// 	}
// 	return true
// }

// func DeleteRunning(c *fiber.Ctx, userid int, level int) error {
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if _, err := DB.QueryContext(ctx, `DELETE FROM running WHERE userid = $1 AND level = $2`, userid, level); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func NewFlag(c *fiber.Ctx, userid int, level int, password string, flag string, port int32, hostname string, deadline int64) error {
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if _, err := DB.QueryContext(ctx, `INSERT INTO flags (userid, level, flag, password, port, hostname, deadline) VALUES ($1, $2, $3, $4, $5, $6, $7)`, userid, level, flag, password, port, hostname, deadline); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func DeleteFlag(c *fiber.Ctx, userid int, level int) error {
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if _, err := DB.QueryContext(ctx, `DELETE FROM flags WHERE userid = $1 AND level = $2`, userid, level); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func ValidChallenge(c *fiber.Ctx, level int) bool {
// 	var chall_name string
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if err := DB.QueryRowContext(ctx, `SELECT chall_name FROM challenges WHERE level = $1`, level).Scan(&chall_name); err != nil {
// 		return false
// 	}
// 	return true
// }

// func ValidFlagEntry(c *fiber.Ctx, level int, userid int) bool {
// 	var flag string
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if err := DB.QueryRowContext(ctx, `SELECT flag FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&flag); err != nil {
// 		return false
// 	}
// 	return true
// }

// func VerifyFlag(c *fiber.Ctx, level int, userid int, flag string) (bool, string) {
// 	var isVerified bool
// 	var acutalflag string
// 	var otheruser int
// 	var currentlevel int
// 	var currentSolves int
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	if err := DB.QueryRowContext(ctx, `SELECT verified FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&isVerified); err != nil {
// 		log.Println(err)
// 		return false, "error in verification, please contact admin"
// 	}
// 	if isVerified {
// 		return false, "flag already verified"
// 	}

// 	if err := DB.QueryRowContext(ctx, `SELECT flag FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&acutalflag); err != nil {
// 		log.Println(err)
// 		return false, "error in verification, please contact admin"
// 	}

// 	if flag == acutalflag {
// 		DB.QueryContext(ctx, `UPDATE flags SET verified = $1 WHERE userid = $2 AND level = $3`, true, userid, level)
// 		if err := DB.QueryRowContext(ctx, `SELECT score FROM users WHERE userid = $1`, userid).Scan(&currentlevel); err != nil {
// 			log.Println(err)
// 			return false, "error in verification, please contact admin"
// 		}
// 		if currentlevel != level {
// 			return false, fmt.Sprintf("Correct flag! no points added. Current level: %d Submitted level: %d", currentlevel, level)
// 		}
// 		DB.QueryContext(ctx, `UPDATE users SET score = $1, lastsubmission = EXTRACT(EPOCH FROM NOW()) WHERE userid = $2`, level+1, userid)

// 		if err := DB.QueryRowContext(ctx, `SELECT solves FROM challenges WHERE level = $1`, level).Scan(&currentSolves); err != nil {
// 			log.Println(err)
// 			return false, "error in verification, please contact admin"
// 		}
// 		DB.QueryContext(ctx, `UPDATE challenges SET solves = $1 WHERE level = $2`, currentSolves+1, level)

// 		return true, "correct flag"
// 	}

// 	if err := DB.QueryRowContext(ctx, `SELECT userid FROM flags WHERE level = $1 AND flag = $2`, level, flag).Scan(&otheruser); err != nil {
// 		return false, "incorrect flag"
// 	}
// 	log.Printf("PLAG: %d submitted %d flag for level %d\n", userid, otheruser, level)
// 	return false, "flag copy detected, incident reported!"
// }

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

// func ReadScores(c *fiber.Ctx) ([]models.Score, error) {
// 	scores := make([]models.Score, 0)
// 	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
// 	defer cancel()

// 	rows, err := DB.QueryContext(ctx, `SELECT username, score from users ORDER BY score DESC, lastsubmission`)
// 	if err != nil {
// 		return scores, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		score := new(models.Score)
// 		if err := rows.Scan(&score.Username, &score.Score); err != nil {
// 			return scores, err
// 		}
// 		scores = append(scores, *score)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return scores, err
// 	}
// 	return scores, nil
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
