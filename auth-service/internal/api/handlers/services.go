package handlers

import (
	"context"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/dto"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/models"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/cache"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/repo"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/services"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/pkg/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UserService interface {
	RegisterUser(ctx context.Context, user dto.UsersRegisterDTO) (*models.Tokens, error)
	LoginUser(ctx context.Context, user dto.UsersLoginDTO) (*models.Tokens, error)
	UpdateAccessToken(ctx context.Context, tokens dto.TokensDTO) (*models.Tokens, error)
	Logout(ctx context.Context, id string) error
}

type Services struct {
	userserv UserService
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
