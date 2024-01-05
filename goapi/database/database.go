package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/CyberLabs-Infosec/isolet/goapi/config"
	"github.com/CyberLabs-Infosec/isolet/goapi/models"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	var err error

	connConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASS, config.DB_NAME)

	DB, err = sql.Open("postgres", connConfig)
	if err != nil {
		return err
	}

	return DB.Ping()
}

func PopulateChalls() error {
	var challenges []models.Challenge

	file, err := os.Open("challenges/challs.json")
	if err != nil {
		return err
	}
	defer file.Close()

	rawData, _ := io.ReadAll(file)
	json.Unmarshal(rawData, &challenges)

	for i := 0; i < len(challenges); i++ {
		if err = AddToChallenges(challenges[i]); err != nil {
			return err
		}
	}
	return nil
}