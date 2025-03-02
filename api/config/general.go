package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var SESSION_SECRET = os.Getenv("SESSION_SECRET")
var ALLOWED_DOMAINS = os.Getenv("ALLOWED_DOMAINS")
var SESSION_EXP = 72
var EMAIL_LEN = 320
var PASS_LEN = 32
var USERNAME_LEN = 32

var REDIS_URL = os.Getenv("REDIS_URL")

var APP_PORT = fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
