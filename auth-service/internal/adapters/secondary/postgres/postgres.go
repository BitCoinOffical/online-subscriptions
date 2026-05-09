package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	timeout = 5
)

func NewPool(cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
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
