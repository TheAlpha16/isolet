package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strconv"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/models"
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

func CheckDomain(email string) bool {
	allowedDomains := []string{"iitism.ac.in"}

	for i := 0; i < len(allowedDomains); i++ {
		domain := allowedDomains[i]
		reg, err := regexp.Compile("^[A-Za-z0-9._%+-]+@" + domain + "$")
		if err != nil {
			log.Println(err)
			return false
		}
		if reg.MatchString(email) {
			return true
		}
	}
	return false
}

func ValidateLoginInput(user *models.User) (bool, string) {
	if len(user.Email) > config.EMAIL_LEN {
		return false, fmt.Sprintf("Email length exceeded %d characters", config.EMAIL_LEN)
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return false, "Not a valid email address"
	}

	if len(user.Password) > config.PASS_LEN {
		return false, "Password length exceeded 32 characters"
	}

	return true, ""
}

func ValidateRegisterInput(regInput *models.User) (bool, string) {
	if len(regInput.Password) > config.PASS_LEN || len(regInput.Password) < 8 {
		return false, fmt.Sprintf("Password should be of 8-%d characters", config.PASS_LEN)
	}

	if regInput.Password != regInput.Confirm {
		return false, "Passwords don't match"
	}

	if len(regInput.Email) > config.EMAIL_LEN {
		return false, fmt.Sprintf("Email length exceeded %d characters", config.EMAIL_LEN)
	}

	if _, err := mail.ParseAddress(regInput.Email); err != nil {
		return false, "Not a valid email address"
	}

	if validDomain := CheckDomain(regInput.Email); !validDomain {
		return false, "Domain is not allowed, please use your iitism.ac.in mail"
	}

	if database.EmailExists(regInput.Email) {
		return false, "Email already exists"
	}

	if len(regInput.Username) > config.USERNAME_LEN {
		return false, fmt.Sprintf("Username exceeded %d characters", config.USERNAME_LEN)
	}

	if len(regInput.Username) < 4 {
		return false, "Username should be of atleast 4 characters"
	}

	if database.UsernameRegistered(regInput.Username, regInput.Email) {
		return false, "Username already exists"
	}

	// if !PassValid(regInput.Password) {
	// 	return false, "Not a strong password"
	// }

	return true, ""
}

func GetInstanceName(userid int, level int) string {
	return Hash(fmt.Sprintf("%d@%d:%s", userid, level, config.INSTANCE_NAME_SECRET))[0:16]
}

func GetHostName(userid int, level int) string {
	return config.INSTANCE_HOSTNAME
}

func BoolAddr(b bool) *bool {
	boolVar := b
	return &boolVar
}

func StringAddr(s string) *string {
	tempString := s
	return &tempString
}

func Int64Addr(i string) *int64 {
	tempInt, _ := strconv.ParseInt(i, 10, 64)
	return &tempInt
}
