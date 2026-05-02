package repo

import (
	"context"
	"fmt"

	"github.com/BitCoinOffical/online-subscriptions/internal/domain/models"
	"github.com/google/uuid"
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
		return uuid.Nil, fmt.Errorf("r.pool.QueryRow: %w", err)
	}
	return id, nil
}
