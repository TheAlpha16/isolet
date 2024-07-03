package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"
	// "github.com/lib/pq"

	"gorm.io/gorm"
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
	var userid int
	err1 := DB.QueryRow(`SELECT userid FROM users WHERE username = $1`, username).Scan(&userid)
	err2 := DB.QueryRow(`SELECT vid FROM toverify WHERE username = $1 AND email != $2`, username, email).Scan(&userid)

	if err1 == nil || err2 == nil {
		return true
	}
	return false
}

func EmailExists(email string) bool {
	var userid int
	if err := DB.QueryRow(`SELECT userid FROM users WHERE email = $1`, email).Scan(&userid); err == nil {
		return true
	}
	return false
}

func AddToVerify(c *fiber.Ctx, user *models.User) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if _, err := DB.QueryContext(ctx, `DELETE FROM toverify WHERE email = $1`, user.Email); err != nil {
		return err
	}
	if _, err := DB.QueryContext(ctx, `INSERT INTO toverify (email, username, password) VALUES ($1, $2, $3)`, user.Email, user.Username, user.Password); err != nil {
		return err
	}
	return nil
}

func AddToUsers(c *fiber.Ctx, email string) (string, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	userData := new(models.User)

	if EmailExists(email) {
		return "user already exists", errors.New("token already verified")
	}

	if err := DB.QueryRowContext(ctx, `SELECT email, username, password FROM toverify WHERE email = $1`, email).Scan(&userData.Email, &userData.Username, &userData.Password); err != nil {
		return "token expired, please register again", err
	}
	if _, err := DB.QueryContext(ctx, `INSERT INTO users (email, username, password) VALUES ($1, $2, $3)`, userData.Email, userData.Username, userData.Password); err != nil {
		log.Println(err.Error())
		return "error in creating user, please contact admin", err
	}
	_, _ = DB.QueryContext(ctx, `DELETE FROM toverify WHERE email = $1`, userData.Email)
	return "", nil
}

func AddToUsersDiscord(c *fiber.Ctx, userid int) (string, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if _, err := DB.QueryContext(ctx, `INSERT INTO users (userid, email, username, password) VALUES ($1, $2, $3, $4)`, userid, strconv.Itoa(userid), strconv.Itoa(userid), GenerateRandom()); err != nil {
		log.Println(err.Error())
		return "error in creating user, please contact admin", err
	}
	return "", nil
}

func AddToChallenges(chall models.Challenge) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := DB.QueryContext(ctx, `INSERT INTO challenges (level, chall_name, prompt, tags) VALUES ($1, $2, $3, $4) ON CONFLICT (level) DO UPDATE SET chall_name = $2, prompt = $3, tags = $4`, chall.Level, chall.Name, chall.Prompt, pq.Array(chall.Tags))
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func ReadChallenges(c *fiber.Ctx) ([]models.Challenge, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	challenges := make([]models.Challenge, 0)
	rows, err := DB.QueryContext(ctx, `SELECT level, chall_name, prompt, solves, tags from challenges ORDER BY level ASC`)
	if err != nil {
		return challenges, err
	}
	defer rows.Close()

	for rows.Next() {
		challenge := new(models.Challenge)
		if err := rows.Scan(&challenge.Level, &challenge.Name, &challenge.Prompt, &challenge.Solves, pq.Array(&challenge.Tags)); err != nil {
			return challenges, err
		}
		challenges = append(challenges, *challenge)
	}
	if err := rows.Err(); err != nil {
		return challenges, err
	}
	return challenges, nil
}

func UserExists(c *fiber.Ctx, userid int) bool {
	var username string
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if err := DB.QueryRowContext(ctx, `SELECT username FROM users WHERE userid = $1`, userid).Scan(&username); err != nil {
		return false
	}
	return true
}

func CanStartInstance(c *fiber.Ctx, userid int, level int) bool {
	var runid int
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if err := DB.QueryRowContext(ctx, `SELECT runid FROM running WHERE userid = $1 AND level = $2`, userid, level).Scan(&runid); err == nil {
		return false
	}

	if _, err := DB.QueryContext(ctx, `INSERT INTO running (userid, level) VALUES ($1, $2)`, userid, level); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func DeleteRunning(c *fiber.Ctx, userid int, level int) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if _, err := DB.QueryContext(ctx, `DELETE FROM running WHERE userid = $1 AND level = $2`, userid, level); err != nil {
		return err
	}
	return nil
}

func NewFlag(c *fiber.Ctx, userid int, level int, password string, flag string, port int32, hostname string, deadline int64) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if _, err := DB.QueryContext(ctx, `INSERT INTO flags (userid, level, flag, password, port, hostname, deadline) VALUES ($1, $2, $3, $4, $5, $6, $7)`, userid, level, flag, password, port, hostname, deadline); err != nil {
		return err
	}
	return nil
}

func DeleteFlag(c *fiber.Ctx, userid int, level int) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if _, err := DB.QueryContext(ctx, `DELETE FROM flags WHERE userid = $1 AND level = $2`, userid, level); err != nil {
		return err
	}
	return nil
}

func ValidChallenge(c *fiber.Ctx, level int) bool {
	var chall_name string
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if err := DB.QueryRowContext(ctx, `SELECT chall_name FROM challenges WHERE level = $1`, level).Scan(&chall_name); err != nil {
		return false
	}
	return true
}

func ValidFlagEntry(c *fiber.Ctx, level int, userid int) bool {
	var flag string
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if err := DB.QueryRowContext(ctx, `SELECT flag FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&flag); err != nil {
		return false
	}
	return true
}

func VerifyFlag(c *fiber.Ctx, level int, userid int, flag string) (bool, string) {
	var isVerified bool
	var acutalflag string
	var otheruser int
	var currentlevel int
	var currentSolves int
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if err := DB.QueryRowContext(ctx, `SELECT verified FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&isVerified); err != nil {
		log.Println(err)
		return false, "error in verification, please contact admin"
	}
	if isVerified {
		return false, "flag already verified"
	}

	if err := DB.QueryRowContext(ctx, `SELECT flag FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&acutalflag); err != nil {
		log.Println(err)
		return false, "error in verification, please contact admin"
	}

	if flag == acutalflag {
		DB.QueryContext(ctx, `UPDATE flags SET verified = $1 WHERE userid = $2 AND level = $3`, true, userid, level)
		if err := DB.QueryRowContext(ctx, `SELECT score FROM users WHERE userid = $1`, userid).Scan(&currentlevel); err != nil {
			log.Println(err)
			return false, "error in verification, please contact admin"
		}
		if currentlevel != level {
			return false, fmt.Sprintf("Correct flag! no points added. Current level: %d Submitted level: %d", currentlevel, level)
		}
		DB.QueryContext(ctx, `UPDATE users SET score = $1, lastsubmission = EXTRACT(EPOCH FROM NOW()) WHERE userid = $2`, level+1, userid)

		if err := DB.QueryRowContext(ctx, `SELECT solves FROM challenges WHERE level = $1`, level).Scan(&currentSolves); err != nil {
			log.Println(err)
			return false, "error in verification, please contact admin"
		}
		DB.QueryContext(ctx, `UPDATE challenges SET solves = $1 WHERE level = $2`, currentSolves+1, level)

		return true, "correct flag"
	}

	if err := DB.QueryRowContext(ctx, `SELECT userid FROM flags WHERE level = $1 AND flag = $2`, level, flag).Scan(&otheruser); err != nil {
		return false, "incorrect flag"
	}
	log.Printf("PLAG: %d submitted %d flag for level %d\n", userid, otheruser, level)
	return false, "flag copy detected, incident reported!"
}

func GetInstances(c *fiber.Ctx, userid int) ([]models.Instance, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	instances := make([]models.Instance, 0)
	rows, err := DB.QueryContext(ctx, `SELECT userid, level, password, port, verified, hostname, deadline from flags WHERE userid = $1`, userid)
	if err != nil {
		return instances, err
	}
	defer rows.Close()

	for rows.Next() {
		instance := new(models.Instance)
		if err := rows.Scan(&instance.UserID, &instance.Level, &instance.Password, &instance.Port, &instance.Verified, &instance.Hostname, &instance.Deadline); err != nil {
			return instances, err
		}
		instances = append(instances, *instance)
	}
	if err := rows.Err(); err != nil {
		return instances, err
	}
	return instances, nil
}

func ReadScores(c *fiber.Ctx) ([]models.Score, error) {
	scores := make([]models.Score, 0)
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	rows, err := DB.QueryContext(ctx, `SELECT username, score from users ORDER BY score DESC, lastsubmission`)
	if err != nil {
		return scores, err
	}
	defer rows.Close()

	for rows.Next() {
		score := new(models.Score)
		if err := rows.Scan(&score.Username, &score.Score); err != nil {
			return scores, err
		}
		scores = append(scores, *score)
	}
	if err := rows.Err(); err != nil {
		return scores, err
	}
	return scores, nil
}

func AddTime(c *fiber.Ctx, userid int, level int) (bool, string, int64) {
	var current int
	var deadline int64
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	if err := DB.QueryRowContext(ctx, `SELECT extended, deadline FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&current, &deadline); err != nil {
		log.Println(err)
		return false, "error in extension, please contact admin", 1
	}

	if (current + 1) > (config.MAX_INSTANCE_TIME / config.INSTANCE_TIME) {
		return false, "limit reached", 1
	}

	_, err := DB.QueryContext(ctx, `UPDATE flags SET extended = $1 WHERE userid = $2 AND level = $3`, current+1, userid, level)
	if err != nil {
		log.Println(err)
		return false, "error in extension, please contact admin", 1
	}

	newdeadline := time.UnixMilli(deadline).Add(time.Minute * time.Duration(config.INSTANCE_TIME)).UnixMilli()

	_, err = DB.QueryContext(ctx, `UPDATE flags SET deadline = $1 WHERE userid = $2 AND level = $3`, newdeadline, userid, level)
	if err != nil {
		log.Println(err)
		return false, "error in extension, please contact admin", 1
	}
	return true, "", newdeadline
}
