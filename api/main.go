package main

import (
	"context"

	"github.com/TheAlpha16/isolet/api/internal/repository"
	"github.com/TheAlpha16/isolet/api/internal/usecase"
	"github.com/TheAlpha16/isolet/api/utils/logger"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	appLogger := logger.GetAppLogger()
	defer appLogger.Sync()

	// Initialize database connection
	dbPool, closeDBConn, err := ConnectToDatabase(ctx)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer closeDBConn()

	// Initialize repositories
	repos := repository.New(dbPool)

	// Initialize usecases
	_ = usecase.New(repos)
}
