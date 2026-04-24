package handlers

import (
	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/repo"
	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Services struct {
	subsserv *services.SubscriptionService
}

func NewServices(pool *pgxpool.Pool) *Services {
	subsrepo := repo.NewSubscriptionHandler(pool)
	subsserv := services.NewSubscriptionService(subsrepo)
	return &Services{subsserv: subsserv}
}

type Handlers struct {
	Subs *SubscriptionHandler
}

func NewHandlers(srvs *Services, logger *zap.Logger) *Handlers {
	subs := NewSubscriptionHandler(srvs.subsserv, logger)
	return &Handlers{Subs: subs}
}
