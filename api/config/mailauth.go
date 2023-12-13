package config

import (
	"os"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
)

var TOKEN_SECRET = os.Getenv("TOKEN_SECRET")
var EMAIL_ID = os.Getenv("EMAIL_ID")
var EMAIL_AUTH = os.Getenv("EMAIL_AUTH")
var TOKEN_EXP = 30
var AUTH_URL = fmt.Sprintf("http://localhost%s/auth/verify?token=", APP_PORT)
var SMTP_HOST = "smtp.gmail.com"
var SMTP_PORT = "587"