package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"avito-task-2025/backend/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	errDbConnect = errors.New("failed to connect to db")
	errDbPing    = errors.New("failed to ping db")
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

	poolConfig, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        return nil, fmt.Errorf("unable to parse config: %w", err)
    }

    poolConfig.MaxConns = int32(cfg.MaxConns)

    pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
    if err != nil {
        return nil, errDbConnect
    }

	err = pool.Ping(ctx)
	if err != nil {
		log.Println(err)
		return nil, errDbPing
	}
	log.Println("success to ping postgres. max conns =", pool.Config().MaxConns)

	return pool, nil
}
