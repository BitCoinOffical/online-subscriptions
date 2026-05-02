package handlers

import (
	"context"

	"github.com/BitCoinOffical/online-subscriptions/internal/domain/dto"
	"github.com/BitCoinOffical/online-subscriptions/internal/domain/models"
	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/cache"
	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/repo"
	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/services"
	"github.com/BitCoinOffical/online-subscriptions/pkg/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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
	userserv *services.UserService
}

func NewServices(pool *pgxpool.Pool, rdb *redis.Client, key string) *Services {
	tokens := jwt.NewManagerToken(key)
	userrepo := repo.NewUserRepo(pool)
	session := cache.NewCache(rdb)
	userserv := services.NewUserService(userrepo, tokens, session)

	subsrepo := repo.NewSubscription(pool)
	subsserv := services.NewSubscriptionService(subsrepo)
	return &Services{subsserv: subsserv, userserv: userserv}
}

type Handlers struct {
	Subs *SubscriptionHandler
	User *User
}

func NewHandlers(srvs *Services, logger *zap.Logger) *Handlers {
	subs := NewSubscriptionHandler(srvs.subsserv, logger)
	user := NewUserHandler(srvs.userserv, logger)

	return &Handlers{Subs: subs, User: user}
}
