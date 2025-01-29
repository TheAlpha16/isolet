package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var SESSION_SECRET = os.Getenv("SESSION_SECRET")
var EVENT_START = os.Getenv("EVENT_START")
var EVENT_END = os.Getenv("EVENT_END")
var POST_EVENT = os.Getenv("POST_EVENT")
var SESSION_EXP = 72
var EMAIL_LEN = 320
var PASS_LEN = 32
var USERNAME_LEN = 32
var TEAM_LEN = 4

var CTF_NAME = os.Getenv("CTF_NAME")
var APP_PORT = fmt.Sprintf(":%s", os.Getenv("APP_PORT"))

var CATEGORY_SIZE = 4
