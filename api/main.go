package main

import (
	"context"
	"fmt"

	restDel "github.com/TheAlpha16/isolet/api/delivery/rest"
	"github.com/TheAlpha16/isolet/api/infra"
	"github.com/TheAlpha16/isolet/api/infra/cache"
	"github.com/TheAlpha16/isolet/api/infra/database/postgres"
	"github.com/TheAlpha16/isolet/api/infra/database/valkey"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/internal/repository"
	"github.com/TheAlpha16/isolet/api/internal/usecase"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	appLogger := logger.GetAppLogger()
	defer appLogger.Sync()

	// Init Sentry
	utils.InitSentry(utils.GetConfig())

	// Initialize database connection
	dbPool, closeDBConn, err := ConnectToPostgresDatabase(ctx)
	if err != nil {
		errorDom.RaiseToSentry(ctx, err)
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer closeDBConn()

	// valkey client
	valkeyClient, err := valkey.NewValkey(ctx)
	if err != nil {
		errorDom.RaiseToSentry(ctx, err)
		appLogger.Fatal("Failed to connect to valkey", zap.Error(err))
	}

	// Init cache
	cache, err := cache.NewClient(ctx, valkeyClient)
	if err != nil {
		errorDom.RaiseToSentry(ctx, err)
		appLogger.Fatal("Failed to connect to cache", zap.Error(err))
	}

	// Init infra services
	infra := infra.New()

	// Init repositories
	repos := repository.New(dbPool)

	// Init usecases
	usecases := usecase.New(cache, repos, infra)

	// Start the rest server
	StartRestServer(ctx, usecases)
}

func ConnectToPostgresDatabase(ctx context.Context) (*gorm.DB, func(), error) {
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
