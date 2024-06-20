package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
	"time"

	"github.com/CyberLabs-Infosec/isolet/goapi/config"
	"github.com/CyberLabs-Infosec/isolet/goapi/models"

	"github.com/golang-jwt/jwt/v5"
)

func getToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Minute * time.Duration(config.TOKEN_EXP)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.TOKEN_SECRET))
	if err != nil {
		return "", err
	}
	return t, nil
}

func SendVerificationMail(user *models.User) error {
	token, err := getToken(user)
	if err != nil {
		return err
	}

	from := config.EMAIL_ID
	secret := config.EMAIL_AUTH

	to := []string{
		user.Email,
	}

	auth := smtp.PlainAuth("", from, secret, config.SMTP_HOST)

	t, _ := template.ParseFiles("templates/mail.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s account verification \n%s\n\n", config.CTF_NAME, mimeHeaders)))

	t.Execute(&body, struct {
		Username string
		Link     string
		Wargame  string
	}{
		Username: user.Username,
		Link:     config.AUTH_URL + token,
		Wargame:  config.CTF_NAME,
	})

	err = smtp.SendMail(config.SMTP_HOST+":"+config.SMTP_PORT, auth, from, to, body.Bytes())
	if err != nil {
		return err
	}

	return nil
}
