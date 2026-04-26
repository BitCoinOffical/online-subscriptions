package services

import (
	"context"

	"github.com/BitCoinOffical/online-subscriptions/internal/domain/models"
)

type SubscriptionRepo interface {
	CreateSubscription(ctx context.Context, sub *models.Subscription) error
	GetSubscriptionsById(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubscriptionsById(ctx context.Context, sub *models.PatchSubscription) error
	FullUpdateSubscriptionsById(ctx context.Context, sub *models.Subscription) error
	DeleteSubscriptions(ctx context.Context, id int) error
	GetSubscriptions(ctx context.Context) ([]models.Subscription, error)
	GetSubscriptionsFilter(ctx context.Context, from, to, user_id, service_name string) (int, error)
}

type Repositories struct {
	repo SubscriptionRepo
}

func NewRepo(subsRepo SubscriptionRepo) *Repositories {
	return &Repositories{repo: subsRepo}
}
