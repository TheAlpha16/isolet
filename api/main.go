package main

import (
	"context"
	"fmt"

	restDel "github.com/TheAlpha16/isolet/api/delivery/rest"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/internal/repository"
	"github.com/TheAlpha16/isolet/api/internal/usecase"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	"github.com/TheAlpha16/isolet/api/utils/postgres"

	"go.uber.org/zap"
	"gorm.io/gorm"
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
	usecases := usecase.New(repos)

	// Start the rest server
	StartRestServer(ctx, usecases)
}

func ConnectToDatabase(ctx context.Context) (*gorm.DB, func(), error) {
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

func StartRestServer(ctx context.Context, usecases *usecase.Usecases) {
	app := restDel.New(usecases)
	restDel.StartServer(app)

	utils.InterruptHandlerChannel <- func() {
		restDel.Shutdown(app)
	}
}
