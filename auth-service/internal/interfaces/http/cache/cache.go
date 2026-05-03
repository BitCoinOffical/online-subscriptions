package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	tokenKey   = "refresh:"
	RefreshTTL = 24 * time.Hour
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{rdb: rdb}
}

func (c *Cache) SaveToken(ctx context.Context, id uuid.UUID, value string) error {
	key := fmt.Sprintf("%s%s", tokenKey, id.String())
	if err := c.rdb.Set(ctx, key, value, RefreshTTL).Err(); err != nil {
		return fmt.Errorf("c.rdb.Set: %w", err)
	}
	return nil
}
