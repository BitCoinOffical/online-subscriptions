package postgres

import (
	"context"
	"fmt"

	"github.com/BitCoinOffical/online-subscriptions/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg *config.PostgresConfig) (*pgxpool.Pool, error) {

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func ClosePool(pool *pgxpool.Pool) {
	pool.Close()
}
