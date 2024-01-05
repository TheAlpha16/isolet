package database

import (
	"database/sql"
	"fmt"

	"github.com/CyberLabs-Infosec/isolet/ripper/config"

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
