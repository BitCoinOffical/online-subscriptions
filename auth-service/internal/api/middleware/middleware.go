package middleware

import (
	"errors"
	"strings"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/api/response"
	jwtpkg "github.com/BitCoinOffical/online-subscriptions/auth-service/pkg/jwt"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	limitergin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
	"go.uber.org/zap"
)

type RateLimiter struct {
	rdb    *redis.Client
	logger *zap.Logger
}

func NewRateLimiter(rdb *redis.Client, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{rdb: rdb, logger: logger}
}

func (r *RateLimiter) RateLimiter() gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted("5-M")
	if err != nil {
		logger.Fatal("limiter.NewRateFromFormatted", zap.Error(err))
	}

	store, err := redisstore.NewStoreWithOptions(r.rdb, limiter.StoreOptions{
		Prefix: "rate_limiter",
	})
	if err != nil {
		logger.Fatal("redisstore.NewStoreWithOptions:", zap.Error(err))
	}
	instance := limiter.New(store, rate)

	return limitergin.NewMiddleware(instance, limitergin.WithLimitReachedHandler(func(c *gin.Context) {
		response.ManyRequest(c, err, "too many requests", r.logger)
		c.Abort()
	}))

}
func AuthMiddleware(jwtManager *jwtpkg.ManagerToken, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(ctx, errors.New("missing token"), "unauthorized", logger)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(ctx, errors.New("invalid token format"), "unauthorized", logger)
			ctx.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(ctx, err, "unauthorized", logger)
			ctx.Abort()
			return
		}
		ctx.Set("user_id", claims.UserID)
		ctx.Next()
	}
}
