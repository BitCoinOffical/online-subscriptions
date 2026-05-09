package handlers

import (
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/cache"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/repo"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/services"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/pkg/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Services struct {
	userserv *services.UserService
}

func NewServices(pool *pgxpool.Pool, rdb *redis.Client, key string) *Services {
	tokens := jwt.NewManagerToken(key)
	userrepo := repo.NewUserRepo(pool)
	session := cache.NewCache(rdb)
	userserv := services.NewUserService(userrepo, tokens, session)

	return &Services{userserv: userserv}
}

type Handlers struct {
	User *User
}

func NewHandlers(srvs *Services, logger *zap.Logger) *Handlers {
	user := NewUserHandler(srvs.userserv, logger)

	return &Handlers{User: user}
}
