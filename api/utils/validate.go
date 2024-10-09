package utils

import (
	"fmt"
	"log"
	"net/mail"
	"regexp"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/models"
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

func ValidateRegisterInput(regInput *models.ToVerify) (bool, string) {
	if len(regInput.Password) > config.PASS_LEN || len(regInput.Password) < 8 {
		return false, fmt.Sprintf("password should be of 8-%d characters", config.PASS_LEN)
	}

	if regInput.Password != regInput.Confirm {
		return false, "passwords don't match"
	}

	if len(regInput.Email) > config.EMAIL_LEN {
		return false, fmt.Sprintf("email length exceeded %d characters", config.EMAIL_LEN)
	}

	if _, err := mail.ParseAddress(regInput.Email); err != nil {
		return false, "not a valid email address"
	}

	// if validDomain := CheckDomain(regInput.Email); !validDomain {
	// 	return false, "Domain is not allowed, please use your iitism.ac.in mail"
	// }

	if database.EmailExists(regInput.Email) {
		return false, "email already exists"
	}

	if len(regInput.Username) > config.USERNAME_LEN {
		return false, fmt.Sprintf("username exceeded %d characters", config.USERNAME_LEN)
	}

	if len(regInput.Username) < 4 {
		return false, "username should be of atleast 4 characters"
	}

	if database.UsernameRegistered(regInput.Username, regInput.Email) {
		return false, "username already exists"
	}

	// if !PassValid(regInput.Password) {
	// 	return false, "Not a strong password"
	// }

	return true, ""
}

