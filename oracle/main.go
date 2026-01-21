package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/TheAlpha16/isolet/oracle/delivery/consumer"
	restDel "github.com/TheAlpha16/isolet/oracle/delivery/rest"
	"github.com/TheAlpha16/isolet/oracle/external"
	"github.com/TheAlpha16/isolet/oracle/infra"
	"github.com/TheAlpha16/isolet/oracle/infra/cache"
	"github.com/TheAlpha16/isolet/oracle/infra/database/postgres"
	"github.com/TheAlpha16/isolet/oracle/infra/database/valkey"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	"github.com/TheAlpha16/isolet/oracle/internal/repository"
	"github.com/TheAlpha16/isolet/oracle/internal/usecase"
	"github.com/TheAlpha16/isolet/oracle/utils"
	"github.com/TheAlpha16/isolet/oracle/utils/logger"
	"github.com/TheAlpha16/isolet/oracle/utils/tracer"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	appLogger := logger.GetAppLogger()
	defer appLogger.Sync()

	var wg sync.WaitGroup

	// Init Sentry
	tracer.InitSentry(utils.GetConfig())

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

	// Start the server based on identity
	switch utils.GetConfig().Identity {
	case utils.IdentityRest:
		StartRestServer(ctx, usecases, infra)
	case utils.IdentityConsumer:
		StartConsumer(ctx, usecases, infra, &wg)
	default:
		appLogger.Fatal("unknown identity", zap.String("identity", string(utils.GetConfig().Identity)))
	}

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

func StartConsumer(ctx context.Context, usecases *usecase.Usecases, infra *infra.Infra, wg *sync.WaitGroup) {
	appLogger := logger.GetAppLogger()
	c, err := consumer.New(ctx, usecases.Fact)
	if err != nil {
		appLogger.Fatal("failed to initialize kafka consumer", zap.Error(err))
	}

	cCtx, cancel := context.WithCancel(ctx)
	utils.InterruptHandlerChannel <- func() {
		cancel()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.Start(cCtx); err != nil {
			errorDom.RaiseToSentry(ctx, err)
			appLogger.Error("kafka consumer failed", zap.Error(err))
		}
	}()
}

/*
TODO
- add metrics to api, herald, tide
- dont init full usecases in case of listener and consumer
*/
