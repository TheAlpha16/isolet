package postgres

import (
	"context"

	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type BaseModel struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}

func NewConnection(ctx context.Context, dbURI string) (*gorm.DB, func(), error) {
	config := utils.GetConfig()

	gormConfig := &gorm.Config{
		Logger:                 glogger.Default.LogMode(glogger.Silent),
		SkipDefaultTransaction: true,
	}
	if config.Environment != utils.PROD {
		gormConfig.Logger = glogger.Default.LogMode(glogger.Info)
	}

	db, err := gorm.Open(postgres.Open(dbURI), gormConfig)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName(config.Database.Name))); err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConnections)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConnections)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(config.Database.MaxConnectionLifeTime)

	err = sqlDB.Ping()
	if err != nil {
		return nil, nil, err
	}
	closeConn := func() {
		sqlDB.Close()
	}

	return db, closeConn, nil
}
