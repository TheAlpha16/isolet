package database

import (
	"fmt"

	"github.com/TheAlpha16/isolet/ripper/config"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error

	connConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASS, config.DB_NAME)

	DB, err = gorm.Open(postgres.Open(connConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	return err
}
