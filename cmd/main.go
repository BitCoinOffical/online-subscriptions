package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BitCoinOffical/online-subscriptions/config"
	"github.com/BitCoinOffical/online-subscriptions/internal/adapters/secondary/migrations"
	"github.com/BitCoinOffical/online-subscriptions/internal/adapters/secondary/postgres"
	"github.com/BitCoinOffical/online-subscriptions/internal/api"
	"github.com/BitCoinOffical/online-subscriptions/internal/api/handlers"
	"github.com/BitCoinOffical/online-subscriptions/pkg"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	migrationsDir = "./migrations"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file: ", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewLoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := pkg.NewLogger(&cfg.App)
	if err != nil {
		log.Fatal(err)
	}
	logger.Debug("data in config", zap.Any("config:", cfg))

	pool, err := postgres.NewPool(ctx, &cfg.Postgres)
	if err != nil {
		logger.Fatal("pool failed", zap.Error(err))
	}
	logger.Info("database pool initialized successfully")

	db := stdlib.OpenDBFromPool(pool)
	if err := migrations.RunMigrations(db, migrationsDir); err != nil {
		logger.Fatal("migrations failed", zap.Error(err))
	}
	logger.Info("database migrations applied successfully")

	srvs := handlers.NewServices(pool)
	handlrs := handlers.NewHandlers(srvs, logger)
	serv := api.NewServer(handlrs)
	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}
}
