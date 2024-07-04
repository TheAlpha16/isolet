package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	// "fmt"
	"log"
	"strconv"
	"time"

	// "github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"
	// "github.com/lib/pq"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GenerateRandom() string {
	buffer := make([]byte, 32)
	rand.Read(buffer)
	return hex.EncodeToString(buffer)
}

func ValidateCreds(c *fiber.Ctx, user *models.User) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Where("(email = ? OR username = ?) AND password = ?", user.Email, user.Email, user.Password).First(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("invalid credentials")
		}
		log.Println(err)
		return err
	}
	return nil
}

func UsernameRegistered(username string, email string) bool {
	var userCount int64
	var toVerifyCount int64

	err1 := DB.Model(&models.User{}).Where("username = ?", username).Count(&userCount).Error
	err2 := DB.Model(&models.ToVerify{}).Where("username = ? AND email != ?", username, email).Count(&toVerifyCount).Error

	if err1 != nil || err2 != nil {
		log.Println(err1, err2)
		return false 
	}

	return userCount > 0 || toVerifyCount > 0
}

func EmailExists(email string) bool {
	var user models.User
	if err := DB.Where("email = ?", email).First(&user).Error; err == nil {
		return true
	}
	return false
}

func AddToVerify(c *fiber.Ctx, toverify *models.ToVerify) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Where("email = ?", toverify.Email).Delete(&models.ToVerify{}).Error; err != nil {
		return err
	}

	if err := db.Create(&toverify).Error; err != nil {
		return err
	}
	return nil
}

func AddToUsers(c *fiber.Ctx, email string) (string, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)
	userData := new(models.User)

	if EmailExists(email) {
		return "user already exists", errors.New("token already verified")
	}

	toVerifyData := new(models.ToVerify)
	if err := db.Where("email = ?", email).First(toVerifyData).Error; err != nil {
		return "token expired, please register again", err
	}

	userData.Email = toVerifyData.Email
	userData.Username = toVerifyData.Username
	userData.Password = toVerifyData.Password

	if err := db.Create(userData).Error; err != nil {
		log.Println(err.Error())
		return "error in creating user, please contact admin", err
	}

	_ = db.Where("email = ?", toVerifyData.Email).Delete(&models.ToVerify{}).Error

	return "", nil
}

func AddToUsersDiscord(c *fiber.Ctx, userid int) (string, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	userData := models.User{
		UserID:   uint(userid),
		Email:    strconv.Itoa(userid),
		Username: strconv.Itoa(userid),
		Password: GenerateRandom(),
	}

	if err := db.Create(&userData).Error; err != nil {
		log.Println(err.Error())
		return "error in creating user, please contact admin", err
	}

	return "", nil
}

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

func ReadChallenges(c *fiber.Ctx) ([]models.Challenge, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var challenges []models.Challenge
	if err := db.Preload("Category").
		Preload("Hints", "visible = ?", true).
		Where("visible = ?", true).
		Find(&challenges).Error; err != nil {
		return challenges, err
	}

	return challenges, nil
}

func UserExists(c *fiber.Ctx, userid int) bool {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var user models.User
	if err := db.Select("username").Where("userid = ?", userid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}
		// Handle other possible errors here if necessary
		return false
	}

	return true
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
