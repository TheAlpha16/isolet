package config

import (
	"os"
	_ "github.com/joho/godotenv/autoload"
)

var DB_USER = os.Getenv("POSTGRES_USER")
var DB_HOST = os.Getenv("POSTGRES_HOST")
var DB_PASS = os.Getenv("POSTGRES_PASSWORD")
var DB_NAME = os.Getenv("POSTGRES_DATABASE")
var DB_PORT = 5432
