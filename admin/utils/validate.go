package utils

import (
	"errors"
	"fmt"

	"github.com/TheAlpha16/isolet/admin/config"
	"github.com/TheAlpha16/isolet/admin/models"
)

// func PassValid(password string) bool {
// 	tests := []string{".{9,}", "[a-z]", "[A-Z]", "[0-9]", "[$#@%&*.!]"}
// 	for _, test := range tests {
// 		t, _ := regexp.MatchString(test, password)
// 		if !t {
// 			return false
// 		}
// 	}
// 	return true
// }

func ValidateLoginInput(user *models.User) (bool, string) {
	if len(user.Email) > config.EMAIL_LEN {
		return false, fmt.Sprintf("email/username length exceeded %d characters", config.EMAIL_LEN)
	}

	// if _, err := mail.ParseAddress(user.Email); err != nil {
	// 	return false, "not a valid email address"
	// }

	if len(user.Password) > config.PASS_LEN {
		return false, "password length exceeded 32 characters"
	}

	return true, ""
}

func ValidateChallengeFields(challenge *models.Challenge) error {
	if challenge.ChallID <= 0 {
		return errors.New("challenge ID is required")
	}

	if challenge.Flag == "" {
		return errors.New("challenge flag cannot be empty")
	}

	if challenge.CategoryID <= 0 || challenge.CategoryID > config.CATEGORY_SIZE {
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

	if challenge.Points <= 0 {
		return errors.New("challenge points must be greater than 0")
	}

	return nil
}



