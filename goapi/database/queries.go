package database

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"strconv"

	// "github.com/CyberLabs-Infosec/isolet/goapi/config"
	"github.com/CyberLabs-Infosec/isolet/goapi/models"
	"github.com/lib/pq"
)

func GenerateRandom() string {
	buffer := make([]byte, 128)
	rand.Read(buffer)
	return fmt.Sprintf("%x", sha256.Sum256(buffer))
}

func ValidateCreds(creds *models.Creds, user *models.User) error {
	if err := DB.QueryRow(`SELECT userid, email, rank FROM users WHERE email = $1 AND password = $2`, creds.Email, creds.Password).Scan(&user.UserID, &user.Email, &user.Rank); err != nil {
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

func AddToVerify(user *models.User) error {
	if _, err := DB.Query(`DELETE FROM toverify WHERE email = $1`, user.Email); err != nil {
		return err
	}
	if _, err := DB.Query(`INSERT INTO toverify (email, username, password) VALUES ($1, $2, $3)`, user.Email, user.Username, user.Password); err != nil {
		return err
	}
	return nil
}

func AddToUsers(email string) (string, error) {
	userData := new(models.User)

	if EmailExists(email) {
		return "user already exists", errors.New("token already verified")
	}

	if err := DB.QueryRow(`SELECT email, username, password FROM toverify WHERE email = $1`, email).Scan(&userData.Email, &userData.Username, &userData.Password); err != nil {
		return "token expired, please register again", err
	}
	if _, err := DB.Query(`INSERT INTO users (email, username, password) VALUES ($1, $2, $3)`, userData.Email, userData.Username, userData.Password); err != nil {
		log.Println(err.Error())
		return "error in creating user, please contact admin", err
	}
	_, _ = DB.Query(`DELETE FROM toverify WHERE email = $1`, userData.Email)
	return "", nil
}

func AddToUsersDiscord(userid int) (string, error) {
	if _, err := DB.Query(`INSERT INTO users (userid, email, username, password) VALUES ($1, $2, $3, $4)`, userid, strconv.Itoa(userid), strconv.Itoa(userid), GenerateRandom()); err != nil {
		log.Println(err.Error())
		return "error in creating user, please contact admin", err
	}
	return "", nil
}

func AddToChallenges(chall models.Challenge) error {
	if _, err := DB.Query(`INSERT INTO challenges (level, chall_name, prompt, tags) VALUES ($1, $2, $3, $4) ON CONFLICT (level) DO UPDATE SET chall_name = $2, prompt = $3, tags = $4`, chall.Level, chall.Name, chall.Prompt, pq.Array(chall.Tags)); err != nil {
		return err
	}
	return nil
}

func ReadChallenges() ([]models.Challenge, error) {
	challenges := make([]models.Challenge, 0)
	rows, err := DB.Query(`SELECT level, chall_name, prompt, solves, tags from challenges ORDER BY level ASC`)
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

func UserExists(userid int) bool {
	var username string
	if err := DB.QueryRow(`SELECT username FROM users WHERE userid = $1`, userid).Scan(&username); err != nil {
		return false
	}
	return true
}

func CanStartInstance(userid int, level int) bool {
	var runid int

	if err := DB.QueryRow(`SELECT runid FROM running WHERE userid = $1 AND level = $2`, userid, level).Scan(&runid); err == nil {
		return false
	}

	if _, err := DB.Query(`INSERT INTO running (userid, level) VALUES ($1, $2)`, userid, level); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func DeleteRunning(userid int, level int) error {
	if _, err := DB.Query(`DELETE FROM running WHERE userid = $1 AND level = $2`, userid, level); err != nil {
		return err
	}
	return nil
}

func NewFlag(userid int, level int, password string, flag string, port int32, hostname string) error {
	if _, err := DB.Query(`INSERT INTO flags (userid, level, flag, password, port, hostname) VALUES ($1, $2, $3, $4, $5, $6)`, userid, level, flag, password, port, hostname); err != nil {
		return err
	}
	return nil
}

func DeleteFlag(userid int, level int) error {
	if _, err := DB.Query(`DELETE FROM flags WHERE userid = $1 AND level = $2`, userid, level); err != nil {
		return err
	}
	return nil
}

func ValidChallenge(level int) bool {
	var chall_name string
	if err := DB.QueryRow(`SELECT chall_name FROM challenges WHERE level = $1`, level).Scan(&chall_name); err != nil {
		return false
	}
	return true
}

func ValidFlagEntry(level int, userid int) bool {
	var flag string
	if err := DB.QueryRow(`SELECT flag FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&flag); err != nil {
		return false
	}
	return true
}

func VerifyFlag(level int, userid int, flag string) (bool, string) {
	var isVerified bool
	var acutalflag string
	var otheruser int
	var currentlevel int
	var currentSolves int

	if err := DB.QueryRow(`SELECT verified FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&isVerified); err != nil {
		log.Println(err)
		return false, "error in verification, please contact admin"
	}
	if isVerified {
		return false, "flag already verified"
	}

	if err := DB.QueryRow(`SELECT flag FROM flags WHERE level = $1 AND userid = $2`, level, userid).Scan(&acutalflag); err != nil {
		log.Println(err)
		return false, "error in verification, please contact admin"
	}

	if flag == acutalflag {
		DB.Query(`UPDATE flags SET verified = $1 WHERE userid = $2 AND level = $3`, true, userid, level)
		if err := DB.QueryRow(`SELECT score FROM users WHERE userid = $1`, userid).Scan(&currentlevel); err != nil {
			log.Println(err)
			return false, "error in verification, please contact admin"
		}
		if currentlevel != level {
			return false, fmt.Sprintf("Correct flag! no points added. Current level: %d Submitted level: %d", currentlevel, level)
		}
		DB.Query(`UPDATE users SET score = $1 WHERE userid = $2`, level+1, userid)

		if err := DB.QueryRow(`SELECT solves FROM challenges WHERE level = $1`, level).Scan(&currentSolves); err != nil {
			log.Println(err)
			return false, "error in verification, please contact admin"
		}
		DB.Query(`UPDATE challenges SET solves = $1 WHERE level = $2`, currentSolves+1, level)

		return true, "correct flag"
	}

	if err := DB.QueryRow(`SELECT userid FROM flags WHERE level = $1 AND flag = $2`, level, flag).Scan(&otheruser); err != nil {
		return false, "incorrect flag"
	}
	log.Printf("PLAG: %d submitted %d flag for level %d\n", userid, otheruser, level)
	return false, "flag copy detected, incident reported!"
}

func GetInstances(userid int) ([]models.Instance, error) {
	instances := make([]models.Instance, 0)
	rows, err := DB.Query(`SELECT userid, level, password, port, verified, hostname from flags WHERE userid = $1`, userid)
	if err != nil {
		return instances, err
	}
	defer rows.Close()

	for rows.Next() {
		instance := new(models.Instance)
		if err := rows.Scan(&instance.UserID, &instance.Level, &instance.Password, &instance.Port, &instance.Verified, &instance.Hostname); err != nil {
			return instances, err
		}
		instances = append(instances, *instance)
	}
	if err := rows.Err(); err != nil {
		return instances, err
	}
	return instances, nil
}

func ReadScores() ([]models.Score, error) {
	scores := make([]models.Score, 0)
	rows, err := DB.Query(`SELECT username, score from users ORDER BY score DESC`)
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