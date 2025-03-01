package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var TOKEN_SECRET = os.Getenv("TOKEN_SECRET")
var TOKEN_EXP = 30
