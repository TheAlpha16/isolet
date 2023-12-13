package utils

import (
	"crypto/sha256"
	"fmt"
	"net/mail"
	"os"
	"regexp"

	"github.com/TitanCrew/isolet/config"
	"github.com/TitanCrew/isolet/database"
	"github.com/TitanCrew/isolet/models"
)

func UpdateKey(key string) error {
	secret := database.GenerateRandom()
	err := os.Setenv(key, secret)
	if err != nil {
		return err
	}
	return nil
}

func Hash(secret string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(secret)))
}

func PassValid(password string) bool {
	tests := []string{".{7,}", "[a-z]", "[A-Z]", "[0-9]", "[$#@%&*.!]"}
	for _, test := range tests {
		t, _ := regexp.MatchString(test, password)
		if !t {
			return false
		}
	}
	return true
}

func ValidateLoginInput(creds *models.Creds) (bool, string) {
	if len(creds.Email) > config.EMAIL_LEN {
		return false, fmt.Sprintf("Email length exceeded %d characters", config.EMAIL_LEN)
	}

	if _, err := mail.ParseAddress(creds.Email); err != nil {
		return false, "Not a valid email address"
	}

	if len(creds.Password) > config.PASS_LEN {
		return false, "Password length exceeded 32 characters"
	}
	
	return true, ""
}

func ValidateRegisterInput(regInput *models.User) (bool, string) {
	if len(regInput.Password) > config.PASS_LEN || len(regInput.Password) < 6 {
		return false, fmt.Sprintf("Password should be of 6-%d characters", config.PASS_LEN)
	}

	if regInput.Password != regInput.Confirm {
		return false, "Passwords doesn't match"
	}

	if len(regInput.Email) > config.EMAIL_LEN {
		return false, fmt.Sprintf("Email length exceeded %d characters", config.EMAIL_LEN)
	}

	if _, err := mail.ParseAddress(regInput.Email); err != nil {
		return false, "Not a valid email address"
	}

	if database.EmailExists(regInput.Email) {
		return false, "Email already exists"
	}

	if len(regInput.Username) > config.USERNAME_LEN {
		return false, fmt.Sprintf("Username exceeded %d characters", config.USERNAME_LEN)
	}

	if len(regInput.Username) < 5 {
		return false, "Username should be of atleast 5 characters"
	}

	if database.UsernameRegistered(regInput.Username, regInput.Email) {
		return false, "Username already exists"
	}

	if !PassValid(regInput.Password) {
		return false, "Not a strong password"
	}

	return true, ""
}

func GetInstanceName(userid int, level int) string {
	return Hash(fmt.Sprintf("%d@%d:%s", userid, level, config.INSTANCE_NAME_SECRET))[0: 16]
}

func BoolAddr(b bool) *bool {
    boolVar := b
    return &boolVar
}

func StringAddr(s string) *string {
	tempString := s
	return &tempString
}