package postgres

import (
	"context"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type BaseModel struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
}

type ImmutableModel struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt int64 `gorm:"autoCreateTime"`
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
		if err := sqlDB.Close(); err != nil {
			errorDom.RaiseToSentry(context.Background(), err)
			logger.GetAppLogger().Error("failed to close PostgreSQL connection", zap.Error(err))
		}
	}

	return db, closeConn, nil
}
