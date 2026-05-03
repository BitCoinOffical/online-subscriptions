package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/config"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/adapters/secondary/migrations"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/adapters/secondary/postgres"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/adapters/secondary/redis"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/api"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/api/handlers"
	zaplogger "github.com/BitCoinOffical/online-subscriptions/auth-service/pkg/logger"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const (
	migrationsDir = "./migrations"
)

// @title Subscriptions API
// @version 1.0
// @description API server for subscriptions

// @host localhost:8080
// @BasePath /api/v1

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewLoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zaplogger.NewLogger(&cfg.App)
	if err != nil {
		log.Fatal(err)
	}
	logger.Debug("data in config", zap.Any("config:", cfg))

	pool, err := postgres.NewPool(&cfg.Postgres)
	if err != nil {
		logger.Fatal("pool failed", zap.Error(err))
	}
	logger.Info("database pool initialized successfully")

	db := stdlib.OpenDBFromPool(pool)
	if err := migrations.RunMigrations(db, migrationsDir); err != nil {
		logger.Fatal("migrations failed", zap.Error(err))
	}
	logger.Info("database migrations applied successfully")

	rdb, err := redis.NewRedis(&cfg.Redis)
	if err != nil {
		logger.Fatal("redis failed", zap.Error(err))
	}
	logger.Info("redis applied successfully")

	srvs := handlers.NewServices(pool, rdb, cfg.App.Jwt)
	handlrs := handlers.NewHandlers(srvs, logger)
	serv := api.NewServer(handlrs, cfg.App.Port)
	go func() {
		if err := serv.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := migrations.RollbackLast(shutdownCtx, db, migrationsDir); err != nil {
		log.Fatalf("goose down failed: %v", err)
	}
	logger.Info("rollback last migrations")

	postgres.ClosePool(pool)
	logger.Info("pool closed")

	if err := serv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}

}
