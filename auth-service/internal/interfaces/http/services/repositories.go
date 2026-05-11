package services

import (
	"context"
	"time"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/models"
	"github.com/google/uuid"
)

type UserRepo interface {
	RegisterUser(ctx context.Context, user *models.Users) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*models.Users, error)
}

type Cache interface {
	SaveToken(ctx context.Context, id uuid.UUID, value string, RefreshTTL time.Duration) error
	GetToken(ctx context.Context, id uuid.UUID) (string, error)
	DeleteRefreshToken(ctx context.Context, id uuid.UUID) error
}
