package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/cache"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	redistc "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	ctx := context.Background()

	container, err := redistc.Run(ctx,
		"redis:7",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections"),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	opt, err := redis.ParseURL(connStr)
	require.NoError(t, err)

	client := redis.NewClient(opt)

	t.Cleanup(func() {
		require.NoError(t, client.Close())
	})

	return client
}

// ──────────────────────────────────────────────
// SaveToken
// ──────────────────────────────────────────────

func TestSaveToken_Success(t *testing.T) {
	client := setupTestRedis(t)
	c := cache.NewCache(client)

	id := uuid.New()
	err := c.SaveToken(context.Background(), id, "refresh-token", time.Minute)

	require.NoError(t, err)

	// проверяем что ключ реально записался
	val := client.Get(context.Background(), "refresh:"+id.String())
	require.NoError(t, val.Err())
	assert.Equal(t, "refresh-token", val.Val())
}

func TestSaveToken_TTLIsSet(t *testing.T) {
	client := setupTestRedis(t)
	c := cache.NewCache(client)

	id := uuid.New()
	ttl := 10 * time.Minute
	err := c.SaveToken(context.Background(), id, "refresh-token", ttl)
	require.NoError(t, err)

	remaining := client.TTL(context.Background(), "refresh:"+id.String())
	require.NoError(t, remaining.Err())
	assert.Greater(t, remaining.Val(), time.Duration(0))
}

// ──────────────────────────────────────────────
// GetToken
// ──────────────────────────────────────────────

func TestGetToken_Success(t *testing.T) {
	client := setupTestRedis(t)
	c := cache.NewCache(client)

	id := uuid.New()
	require.NoError(t, c.SaveToken(context.Background(), id, "refresh-token", time.Minute))

	token, err := c.GetToken(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, "refresh-token", token)
}

func TestGetToken_NotFound(t *testing.T) {
	client := setupTestRedis(t)
	c := cache.NewCache(client)

	token, err := c.GetToken(context.Background(), uuid.New())

	assert.Empty(t, token)
	assert.ErrorContains(t, err, "c.rdb.Get")
}

func TestGetToken_AfterExpiry(t *testing.T) {
	client := setupTestRedis(t)
	c := cache.NewCache(client)

	id := uuid.New()
	require.NoError(t, c.SaveToken(context.Background(), id, "refresh-token", time.Millisecond*100))

	time.Sleep(time.Millisecond * 200)

	token, err := c.GetToken(context.Background(), id)

	assert.Empty(t, token)
	assert.ErrorContains(t, err, "c.rdb.Get")
}

// ──────────────────────────────────────────────
// DeleteRefreshToken
// ──────────────────────────────────────────────

func TestDeleteRefreshToken_Success(t *testing.T) {
	client := setupTestRedis(t)
	c := cache.NewCache(client)

	id := uuid.New()
	require.NoError(t, c.SaveToken(context.Background(), id, "refresh-token", time.Minute))

	err := c.DeleteRefreshToken(context.Background(), id)
	require.NoError(t, err)

	// ключ должен исчезнуть
	val := client.Get(context.Background(), "refresh:"+id.String())
	assert.ErrorIs(t, val.Err(), redis.Nil)
}

func TestDeleteRefreshToken_NonExistentKey(t *testing.T) {
	client := setupTestRedis(t)
	c := cache.NewCache(client)

	// Redis не возвращает ошибку при DEL несуществующего ключа — это нормальное поведение
	err := c.DeleteRefreshToken(context.Background(), uuid.New())

	assert.NoError(t, err)
}
