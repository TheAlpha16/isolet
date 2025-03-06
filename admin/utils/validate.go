package utils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/TheAlpha16/isolet/admin/config"
	"github.com/TheAlpha16/isolet/admin/models"
)

func ValidateLoginInput(user *models.User) (bool, string) {
	if len(user.Email) > config.EMAIL_LEN {
		return false, fmt.Sprintf("email/username length exceeded %d characters", config.EMAIL_LEN)
	}

	if len(user.Password) > config.PASS_LEN {
		return false, "password length exceeded 32 characters"
	}

	return true, ""
}

func ValidateChallengeFields(challenge *models.Challenge) error {
	if challenge.ChallID <= 0 && reflect.TypeOf(challenge.ChallID).Kind() == reflect.Int{
		return errors.New("challenge ID is required")
	}

	if (challenge.CategoryID <= 0 || challenge.CategoryID > config.CATEGORY_SIZE) && reflect.TypeOf(challenge.CategoryID).Kind() == reflect.Int{
		return errors.New("category id out of bounds")
	}

	if challenge.Type != "" {
		validTypes := map[string]bool{
			"static":    true,
			"dynamic":   true,
			"on-demand": true,
		}
		if !validTypes[challenge.Type] {
			return errors.New("invalid challenge type")
		}
	}

	if challenge.Points <= 0 && reflect.TypeOf(challenge.Points).Kind() == reflect.Int{
		return errors.New("challenge points must be greater than 0")
	}

	return nil
}

func ValidateChallengeFileFields(challenge *models.Challenge) error {
	if challenge.ChallID <= 0 && reflect.TypeOf(challenge.ChallID).Kind() == reflect.Int{
		return errors.New("challenge ID is required")
	}

	if (challenge.CategoryID <= 0 || challenge.CategoryID > config.CATEGORY_SIZE) && reflect.TypeOf(challenge.CategoryID).Kind() == reflect.Int{
		return errors.New("category id out of bounds")
	}

	if challenge.Type != "" {
		validTypes := map[string]bool{
			"static":    true,
			"dynamic":   true,
			"on-demand": true,
		}
		if !validTypes[challenge.Type] {
			return errors.New("invalid challenge type")
		}
	}

	return nil
}

func ValidateConfigFields (config *models.Config) error {
	if config.Key == "" {
		return errors.New("config key cannot be empty")
	}

	if config.Value == "" {
		return errors.New("config value cannot be empty")
	}

	return nil
}

func ValidateHintFields (hint *models.Hint) error {
	if hint.HID <= 0 && reflect.TypeOf(hint.HID).Kind() == reflect.Int{
		return errors.New("HID is required")
	}

	if hint.Cost < 0 {
		return errors.New("hint cost cannot be less than 0")
	}

	return nil
}
