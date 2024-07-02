package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"

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

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return DB.PingContext(ctx)
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
