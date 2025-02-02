package database

import (
	"context"
	"errors"

	// "fmt"
	"log"
	"time"

	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

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
		log.Println(err)
		return "error in creating user, please contact admin", err
	}

	_ = db.Where("email = ?", toVerifyData.Email).Delete(&models.ToVerify{}).Error

	return "", nil
}

func UserExists(c *fiber.Ctx, userid int64) bool {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var user models.User
	if err := db.Select("username").Where("userid = ?", userid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false
		}

		return false
	}

	return true
}

func UserInTeam(c *fiber.Ctx, userid int64) bool {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var user models.User
	if err := db.Select("teamid").Where("userid = ?", userid).First(&user).Error; err != nil {
		return err == gorm.ErrRecordNotFound
	}

	return user.TeamID != -1
}

func UpdateUserTeam(c *fiber.Ctx, userid int64, teamid int64) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Model(&models.User{}).Where("userid = ?", userid).Update("teamid", teamid).Error; err != nil {
		return err
	}

	return nil
}

func GeneratePasswordResetToken(c *fiber.Ctx, user *models.User) (*models.Token, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)
	var token models.Token

	if err := db.Where("email = ?", user.Email).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Println(err)
		}
		return &token, err
	}

	token.UserID = user.UserID
	token.Type = "password_reset"
	token.Expiry = time.Now().Add(30 * time.Minute)

	if err := db.Create(&token).Error; err != nil {
		log.Println(err)
		return &token, err
	}

	return &token, nil
}

func VerifyToken(c *fiber.Ctx, token string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var tokenData models.Token

	if err := db.Preload("User").Where("token = ?", token).First(&tokenData).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Println(err)
		}
		return &tokenData.User, errors.New("invalid token")
	}

	if tokenData.Type != "password_reset" {
		return &tokenData.User, errors.New("invalid token")
	}

	if tokenData.Expiry.Before(time.Now()) {
		return &tokenData.User, errors.New("token expired")
	}

	if err := db.Delete(&tokenData).Error; err != nil {
		log.Println(err)
		return &tokenData.User, errors.New("invalid token")
	}

	return &tokenData.User, nil
}

func ResetPassword(c *fiber.Ctx, user *models.User) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	if err := db.Model(&models.User{}).Where("email = ?", user.Email).Update("password", user.Password).Error; err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func ResetForgetPassword(c *fiber.Ctx, token string, password string) error {
	user, err := VerifyToken(c, token)
	if err != nil {
		return err
	}

	user.Password = password

	if err := ResetPassword(c, user); err != nil {
		return errors.New("error in resetting password")
	}

	return nil
}