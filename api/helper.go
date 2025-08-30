package main

import (
	"context"
	"fmt"

	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
