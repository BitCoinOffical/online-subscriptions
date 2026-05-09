package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain/models"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/interfaces/http/repo"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:18",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS subscriptions (
			id BIGSERIAL PRIMARY KEY,
			service_name VARCHAR(255) NOT NULL,
			price INTEGER NOT NULL,
			user_id UUID NOT NULL,
			start_date DATE NOT NULL,
			end_date DATE,

			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
    `)
	require.NoError(t, err)

	return pool
}

func TestCreateSubscription(t *testing.T) {
	pool := setupTestDB(t)
	repo := repo.NewSubscription(pool)

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			pool.Exec(context.Background(), "TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE")
		})
		sub := &models.Subscription{
			ServiceName: "Netflix",
			Price:       100,
			UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		}

		err := repo.CreateSubscription(context.Background(), sub)
		assert.NoError(t, err)

	})
}

func TestGetSubscriptionsById(t *testing.T) {
	pool := setupTestDB(t)
	r := repo.NewSubscription(pool)

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			pool.Exec(context.Background(), "TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE")
		})
		endDate := time.Now().AddDate(1, 0, 0)
		sub := &models.Subscription{
			ServiceName: "Netflix",
			Price:       100,
			UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			StartDate:   time.Now(),
			EndDate:     &endDate,
		}
		err := r.CreateSubscription(context.Background(), sub)
		require.NoError(t, err)

		got, err := r.GetSubscriptionsById(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, "Netflix", got.ServiceName)
	})

	t.Run("not found", func(t *testing.T) {
		got, err := r.GetSubscriptionsById(context.Background(), 999)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, got)
	})
}

func TestUpdateSubscriptionsById(t *testing.T) {
	pool := setupTestDB(t)
	r := repo.NewSubscription(pool)

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			pool.Exec(context.Background(), "TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE")
		})
		endDate := time.Now().AddDate(1, 0, 0)
		sub := &models.Subscription{
			ServiceName: "Netflix",
			Price:       100,
			UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			StartDate:   time.Now(),
			EndDate:     &endDate,
		}
		err := r.CreateSubscription(context.Background(), sub)
		require.NoError(t, err)
		price := 10
		err = r.UpdateSubscriptionsById(context.Background(), &models.PatchSubscription{
			ID:    1,
			Price: &price,
		})
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		price := 10
		err := r.UpdateSubscriptionsById(context.Background(), &models.PatchSubscription{
			ID:    999,
			Price: &price,
		})
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestFullUpdateSubscriptionsById(t *testing.T) {
	pool := setupTestDB(t)
	r := repo.NewSubscription(pool)

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			pool.Exec(context.Background(), "TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE")
		})
		endDate := time.Now().AddDate(1, 0, 0)
		sub := &models.Subscription{
			ServiceName: "Netflix",
			Price:       100,
			UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			StartDate:   time.Now(),
			EndDate:     &endDate,
		}
		err := r.CreateSubscription(context.Background(), sub)
		require.NoError(t, err)

		err = r.FullUpdateSubscriptionsById(context.Background(), &models.Subscription{
			ID:          1,
			ServiceName: "HBO",
			Price:       200,
			UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			StartDate:   time.Now(),
			EndDate:     &endDate,
		})
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		endDate := time.Now().AddDate(1, 0, 0)
		err := r.FullUpdateSubscriptionsById(context.Background(), &models.Subscription{
			ID:          999,
			ServiceName: "HBO",
			Price:       200,
			UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			StartDate:   time.Now(),
			EndDate:     &endDate,
		})
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestDeleteSubscriptions(t *testing.T) {
	pool := setupTestDB(t)
	r := repo.NewSubscription(pool)

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			pool.Exec(context.Background(), "TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE")
		})
		endDate := time.Now().AddDate(1, 0, 0)
		sub := &models.Subscription{
			ServiceName: "Netflix",
			Price:       100,
			UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			StartDate:   time.Now(),
			EndDate:     &endDate,
		}
		err := r.CreateSubscription(context.Background(), sub)
		require.NoError(t, err)

		err = r.DeleteSubscriptions(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		err := r.DeleteSubscriptions(context.Background(), 999)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestGetSubscriptions(t *testing.T) {
	pool := setupTestDB(t)
	r := repo.NewSubscription(pool)

	t.Run("success", func(t *testing.T) {
		t.Cleanup(func() {
			pool.Exec(context.Background(), "TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE")
		})
		endDate := time.Now().AddDate(1, 0, 0)
		sub := &models.Subscription{
			ServiceName: "Netflix",
			Price:       100,
			UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			StartDate:   time.Now(),
			EndDate:     &endDate,
		}
		err := r.CreateSubscription(context.Background(), sub)
		require.NoError(t, err)

		subs, err := r.GetSubscriptions(context.Background())
		assert.NoError(t, err)
		assert.Len(t, subs, 1)
		assert.Equal(t, "Netflix", subs[0].ServiceName)
	})

	t.Run("empty", func(t *testing.T) {
		subs, err := r.GetSubscriptions(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, subs)
	})
}

func TestGetSubscriptionsFilter(t *testing.T) {
	pool := setupTestDB(t)
	r := repo.NewSubscription(pool)

	t.Cleanup(func() {
		pool.Exec(context.Background(), "TRUNCATE TABLE subscriptions RESTART IDENTITY CASCADE")
	})

	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	endDate := time.Now().AddDate(1, 0, 0)
	sub := &models.Subscription{
		ServiceName: "Netflix",
		Price:       100,
		UserID:      userID,
		StartDate:   time.Now(),
		EndDate:     &endDate,
	}
	err := r.CreateSubscription(context.Background(), sub)
	require.NoError(t, err)

	t.Run("success all filters", func(t *testing.T) {
		total, err := r.GetSubscriptionsFilter(context.Background(),
			time.Now().AddDate(0, -1, 0).Format("2006-01-02"),
			time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
			userID.String(),
			"Netflix",
		)
		assert.NoError(t, err)
		assert.Equal(t, 100, total)
	})

	t.Run("no filters", func(t *testing.T) {
		total, err := r.GetSubscriptionsFilter(context.Background(), "", "", "", "")
		assert.NoError(t, err)
		assert.Equal(t, 100, total)
	})

	t.Run("no match", func(t *testing.T) {
		total, err := r.GetSubscriptionsFilter(context.Background(), "", "", "", "HBO")
		assert.NoError(t, err)
		assert.Equal(t, 0, total)
	})
}
