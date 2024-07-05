package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/TheAlpha16/isolet/ripper/config"

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

	DB.SetMaxOpenConns(5)
	DB.SetMaxIdleConns(5)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return DB.PingContext(ctx)
}
