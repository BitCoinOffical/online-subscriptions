package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) RegisterUser(ctx context.Context, user *models.Users) (uuid.UUID, error) {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query, user.Email, user.Password_hash).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return uuid.Nil, fmt.Errorf("email already exists: %w", domain.ErrEmailAlreadyExists)
		}
		return uuid.Nil, fmt.Errorf("r.pool.QueryRow: %w", err)
	}
	return id, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.Users, error) {
	query := `SELECT id, password_hash FROM users WHERE email = $1`

	var user models.Users

	err := r.pool.QueryRow(ctx, query, email).Scan(&user.Id, &user.Password_hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("r.pool.QueryRow: %w", err)
	}

	return &user, nil
}
