package main

import (
	"context"
	"fmt"
	"sync"

	restDel "github.com/TheAlpha16/isolet/api/delivery/rest"
	"github.com/TheAlpha16/isolet/api/external"
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

	var wg sync.WaitGroup

	// Init Sentry
	utils.InitSentry(utils.GetConfig())

	// Initialize database connection
	dbPool, closeDBConn, err := ConnectToPostgresDatabase(ctx)
	if err != nil {
		errorDom.RaiseToSentry(ctx, err)
		appLogger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer closeDBConn()

	// valkey client
	valkeyClient, err := valkey.NewValkey(ctx)
	if err != nil {
		errorDom.RaiseToSentry(ctx, err)
		appLogger.Fatal("failed to connect to valkey", zap.Error(err))
	}

	// Init cache
	cache, err := cache.NewClient(ctx, valkeyClient)
	if err != nil {
		errorDom.RaiseToSentry(ctx, err)
		appLogger.Fatal("failed to connect to cache", zap.Error(err))
	}

	// automigrate database
	if utils.GetConfig().Environment == utils.LOCAL {
		err = repository.AutoMigrate(dbPool)
		if err != nil {
			errorDom.RaiseToSentry(ctx, err)
			appLogger.Fatal("failed to migrate database", zap.Error(err))
		}
	}

	// Init infra services
	infra, err := infra.New(ctx, valkeyClient)
	if err != nil {
		errorDom.RaiseToSentry(ctx, err)
		appLogger.Fatal("failed to initialize infra services", zap.Error(err))
	}

	// External services
	external, err := external.New(ctx)
	if err != nil {
		errorDom.RaiseToSentry(ctx, err)
		appLogger.Fatal("failed to initialize external services", zap.Error(err))
	}

	// Init repositories
	repos := repository.New(dbPool, cache)

	// Init usecases
	usecases := usecase.New(ctx, &wg, cache, repos, infra, external)

	// Start the rest server
	StartRestServer(ctx, usecases, infra)

	utils.InterruptHandlerChannel <- func() {
		wg.Wait()
	}
	utils.InterruptHandler()
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

func StartRestServer(ctx context.Context, usecases *usecase.Usecases, infra *infra.Infra) {
	app := restDel.New(usecases, infra)
	utils.InterruptHandlerChannel <- func() {
		restDel.Shutdown(app)
	}
	go restDel.StartServer(app)
}
