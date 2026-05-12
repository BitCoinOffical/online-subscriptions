package repo_test

import (
	"context"
	"testing"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/models"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
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
		require.NoError(t, container.Terminate(ctx))
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email         VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL
		)
	`)
	require.NoError(t, err)

	return pool
}

func cleanUsers(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `TRUNCATE TABLE users`)
	require.NoError(t, err)
}

func TestRegisterUser_Success(t *testing.T) {
	pool := setupTestDB(t)
	r := repo.NewUserRepo(pool)

	user := &models.Users{
		Email:         "test@example.com",
		Password_hash: "hashed_password",
	}

	id, err := r.RegisterUser(context.Background(), user)

	require.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	pool := setupTestDB(t)
	r := repo.NewUserRepo(pool)

	user := &models.Users{
		Email:         "duplicate@example.com",
		Password_hash: "hashed_password",
	}

	_, err := r.RegisterUser(context.Background(), user)
	require.NoError(t, err)

	_, err = r.RegisterUser(context.Background(), user)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrEmailAlreadyExists)
}

func TestGetUserByEmail_Success(t *testing.T) {
	pool := setupTestDB(t)
	r := repo.NewUserRepo(pool)

	user := &models.Users{
		Email:         "found@example.com",
		Password_hash: "hashed_password",
	}

	expectedID, err := r.RegisterUser(context.Background(), user)
	require.NoError(t, err)

	result, err := r.GetUserByEmail(context.Background(), user.Email)

	require.NoError(t, err)
	assert.Equal(t, expectedID, result.Id)
	assert.Equal(t, user.Password_hash, result.Password_hash)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	r := repo.NewUserRepo(pool)

	result, err := r.GetUserByEmail(context.Background(), "ghost@example.com")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
