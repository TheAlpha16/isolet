package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"

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

func SendVerificationMail(regInput *models.ToVerify) error {
	user := &models.User{
		Email: regInput.Email,
	}

	token, err := getToken(user)
	if err != nil {
		return err
	}

	ctfName, err := config.Get("CTF_NAME")
	if err != nil {
		return err
	}

	emailID, err := config.Get("EMAIL_ID")
	if err != nil {
		return err
	}

	emailAuth, err := config.Get("EMAIL_AUTH")
	if err != nil {
		return err
	}

	smtpHost, err := config.Get("SMTP_HOST")
	if err != nil {
		return err
	}

	smtpPort, err := config.Get("SMTP_PORT")
	if err != nil {
		return err
	}

	publicURL, err := config.Get("PUBLIC_URL")
	if err != nil {
		return err
	}

	emailUsername, err := config.Get("EMAIL_USERNAME")
	if err != nil {
		return err
	}

	from := emailID
	secret := emailAuth

	to := []string{
		regInput.Email,
	}

	auth := smtp.PlainAuth("", emailUsername, secret, smtpHost)

	t, err := template.ParseFiles("templates/mail.html")
	if err != nil {
		log.Println(err)
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s account verification \n%s\n\n", ctfName, mimeHeaders)))

	err = t.Execute(&body, struct {
		Username string
		Link     string
		Wargame  string
	}{
		Username: regInput.Username,
		Link:     fmt.Sprintf("http://%s/auth/verify?token=%s", publicURL, token),
		Wargame:  ctfName,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func SendResetPasswordMail(user *models.User, token *models.Token) error {
	var err error

	emailID, err := config.Get("EMAIL_ID")
	if err != nil {
		return err
	}

	emailAuth, err := config.Get("EMAIL_AUTH")
	if err != nil {
		return err
	}

	smtpHost, err := config.Get("SMTP_HOST")
	if err != nil {
		return err
	}

	smtpPort, err := config.Get("SMTP_PORT")
	if err != nil {
		return err
	}

	publicURL, err := config.Get("PUBLIC_URL")
	if err != nil {
		return err
	}

	emailUsername, err := config.Get("EMAIL_USERNAME")
	if err != nil {
		return err
	}

	from := emailID
	secret := emailAuth
	ctfName, err := config.Get("CTF_NAME")
	if err != nil {
		log.Println(err)
		return err
	}

	to := []string{
		user.Email,
	}

	auth := smtp.PlainAuth("", emailUsername, secret, smtpHost)

	t, err := template.ParseFiles("templates/reset.html")
	if err != nil {
		log.Println(err)
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s password reset \n%s\n\n", ctfName, mimeHeaders)))

	err = t.Execute(&body, struct {
		Link     string
		Username string
	}{
		Username: user.Username,
		Link:     fmt.Sprintf("http://%s/reset-password?token=%s", publicURL, token.Token),
	})
	if err != nil {
		log.Println(err)
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
