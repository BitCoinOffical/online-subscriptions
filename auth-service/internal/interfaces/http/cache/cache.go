package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	tokenKey = "refresh:"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{rdb: rdb}
}

func (c *Cache) SaveToken(ctx context.Context, id uuid.UUID, value string, RefreshTTL time.Duration) error {
	key := fmt.Sprintf("%s%s", tokenKey, id.String())
	if err := c.rdb.Set(ctx, key, value, RefreshTTL).Err(); err != nil {
		return fmt.Errorf("c.rdb.Set: %w", err)
	}
	return nil
}
func (c *Cache) GetToken(ctx context.Context, id uuid.UUID) (string, error) {
	key := fmt.Sprintf("%s%s", tokenKey, id.String())
	token, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("c.rdb.Get: %w", err)
	}
	return token, nil
}
func (c *Cache) DeleteRefreshToken(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("%s%s", tokenKey, id.String())
	if err := c.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("c.rdb.Del: %w", err)
	}
	return nil
}
