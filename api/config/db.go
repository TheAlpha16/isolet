package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var DB_USER = os.Getenv("POSTGRES_USER")
var DB_HOST = os.Getenv("POSTGRES_HOST")
var DB_PASS = os.Getenv("POSTGRES_PASSWORD")
var DB_NAME = os.Getenv("POSTGRES_DATABASE")
var DB_PORT = getPort()

func getPort() int {
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return 5432
	}
	return port
}