package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/BitCoinOffical/online-subscriptions/subscription-service/docs"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/config"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/adapters/secondary/migrations"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/adapters/secondary/postgres"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/adapters/secondary/redis"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/api"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/api/handlers"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/api/middleware"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/pkg/jwt"
	zaplogger "github.com/BitCoinOffical/online-subscriptions/subscription-service/pkg/logger"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const (
	migrationsDir = "./migrations"
)

// @title Subscription Service API
// @version 1.0
// @description API server for managing user subscriptions
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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

	manager := jwt.NewManagerToken(cfg.App.Jwt)

	srvs := handlers.NewServices(pool)
	handlrs := handlers.NewHandlers(srvs, logger)
	limiter := middleware.NewRateLimiter(rdb, logger)
	serv := api.NewServer(handlrs, cfg.App.Port, manager, limiter, logger)
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
