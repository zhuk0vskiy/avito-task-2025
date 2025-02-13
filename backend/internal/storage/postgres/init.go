package postgres

import (
	"context"
	"errors"
	"fmt"

	"avito-task-2025/backend/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errDbConnect = errors.New("failed to connect to db")
	errDbPing = errors.New("failed to ping db")
)

func NewDbConn(ctx context.Context, cfg *config.PostgresConfig) (pool *pgxpool.Pool, err error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		cfg.Driver,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	pool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, errDbConnect
	}

	err = pool.Ping(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, errDbPing
	}

	return pool, nil
}