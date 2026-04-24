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
	"github.com/BitCoinOffical/online-subscriptions/internal/rules"
	"github.com/BitCoinOffical/online-subscriptions/pkg"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const (
	migrationsDir = "./migrations"
)

func main() {

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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("price", rules.ValidatePrice); err != nil {
			logger.Error("failed to register price validator", zap.Error(err))
		}
	}

	srvs := handlers.NewServices(pool)
	handlrs := handlers.NewHandlers(srvs, logger)
	serv := api.NewServer(handlrs)
	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}
}
