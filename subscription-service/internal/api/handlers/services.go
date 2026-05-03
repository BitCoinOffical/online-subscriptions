package handlers

import (
	"context"

	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain/dto"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain/models"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/interfaces/http/repo"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/interfaces/http/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, dto *dto.SubscriptionDTO) error
	GetSubscriptionsById(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubscriptionsById(ctx context.Context, dto *dto.PatchSubscriptionDTO, id int) error
	FullUpdateSubscriptionsById(ctx context.Context, dto *dto.SubscriptionDTO, id int) error
	DeleteSubscriptions(ctx context.Context, id int) error
	GetSubscriptions(ctx context.Context) ([]models.Subscription, error)
	GetSubscriptionsFilter(ctx context.Context, from, to, user_id, service_name string) (int, error)
}

type Services struct {
	subsserv SubscriptionService
}

func NewServices(pool *pgxpool.Pool) *Services {
	subsrepo := repo.NewSubscription(pool)
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
