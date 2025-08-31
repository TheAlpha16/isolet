package main

import (
	"context"
	"fmt"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/internal/repository"
	"github.com/TheAlpha16/isolet/api/internal/usecase"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	"github.com/TheAlpha16/isolet/api/utils/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	appLogger := logger.GetAppLogger()
	defer appLogger.Sync()

	// Initialize Sentry
	utils.InitSentry(utils.GetConfig())

	// Initialize database connection
	dbPool, closeDBConn, err := ConnectToDatabase(ctx)
	if err != nil {
		errorDom.RaiseToSentry(ctx, err)
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer closeDBConn()

	// Initialize repositories
	repos := repository.New(dbPool)

	// Initialize usecases
	_ = usecase.New(repos)
}

func ConnectToDatabase(ctx context.Context) (*pgxpool.Pool, func(), error) {
	config := utils.GetConfig()
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	dbPool, closeConn, err := postgres.NewConnection(ctx, dbURI)
	if err != nil {
		return nil, nil, err
	}

	return dbPool, closeConn, nil
}
