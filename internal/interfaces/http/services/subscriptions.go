package services

import (
	"context"
	"fmt"

	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/models"
	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/repo"
)

type SubscriptionService struct {
	repo *repo.SubscriptionRepo
}

func NewSubscriptionService(repo *repo.SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, dto *dto.SubscriptionDTO) error {
	sub := models.Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      dto.UserID,
		StartDate:   dto.StartDate,
		EndDate:     dto.EndDate,
	}
	err := s.repo.CreateSubscription(ctx, &sub)
	if err != nil {
		return fmt.Errorf("s.repo.CreateSubscription: %w", err)
	}
	return nil
}

func (s *SubscriptionService) GetSubscriptionsById(ctx context.Context, id int) (*models.Subscription, error) {
	sub, err := s.repo.GetSubscriptionsById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetSubscriptionsById: %w", err)
	}
	return sub, nil
}

func (s *SubscriptionService) UpdateSubscriptions(ctx context.Context, dto *dto.SubscriptionDTO, id int) error {
	sub := models.Subscription{
		ID:          id,
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      dto.UserID,
		StartDate:   dto.StartDate,
		EndDate:     dto.EndDate,
	}
	err := s.repo.UpdateSubscriptions(ctx, &sub)
	if err != nil {
		return fmt.Errorf("s.repo.UpdateSubscriptions: %w", err)
	}
	return nil
}

func (s *SubscriptionService) DeleteSubscriptions(ctx context.Context, id int) error {
	err := s.repo.DeleteSubscriptions(ctx, id)
	if err != nil {
		return fmt.Errorf("s.repo.DeleteSubscriptions: %w", err)
	}
	return nil
}

func (s *SubscriptionService) GetSubscriptions(ctx context.Context) ([]models.Subscription, error) {
	subs, err := s.repo.GetSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.repo.GetSubscriptions: %w", err)
	}
	return subs, nil
}
