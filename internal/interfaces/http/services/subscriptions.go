package services

import (
	"context"
	"fmt"

	"github.com/BitCoinOffical/online-subscriptions/internal/domain/dto"
	"github.com/BitCoinOffical/online-subscriptions/internal/domain/models"
	"github.com/BitCoinOffical/online-subscriptions/pkg/date"
)

type SubscriptionService struct {
	repo SubscriptionRepo
}

func NewSubscriptionService(repo SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, dto *dto.SubscriptionDTO) error {
	startT, err := date.ParseMonthDate(dto.StartDate)
	if err != nil {
		return err
	}

	endT, err := date.ParseMonthDate(dto.EndDate)
	if err != nil {
		return err
	}

	sub := models.Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      dto.UserID,
		StartDate:   startT,
		EndDate:     endT,
	}
	err = s.repo.CreateSubscription(ctx, &sub)
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

func (s *SubscriptionService) UpdateSubscriptionsById(ctx context.Context, dto *dto.PatchSubscriptionDTO, id int) error {

	startT, err := date.ParseMonthDate(dto.StartDate)
	if err != nil {
		return fmt.Errorf("date.ParseMonthDate: %w", err)
	}

	endT, err := date.ParseMonthDate(dto.EndDate)
	if err != nil {
		return fmt.Errorf("date.ParseMonthDate: %w", err)
	}

	sub := models.PatchSubscription{
		ID:          id,
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      dto.UserID,
		StartDate:   startT,
		EndDate:     endT,
	}
	err = s.repo.UpdateSubscriptionsById(ctx, &sub)
	if err != nil {
		return fmt.Errorf("s.repo.UpdateSubscriptions: %w", err)
	}
	return nil
}

func (s *SubscriptionService) FullUpdateSubscriptionsById(ctx context.Context, dto *dto.SubscriptionDTO, id int) error {
	startT, err := date.ParseMonthDate(dto.StartDate)
	if err != nil {
		return err
	}

	endT, err := date.ParseMonthDate(dto.EndDate)
	if err != nil {
		return err
	}

	sub := models.Subscription{
		ID:          id,
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      dto.UserID,
		StartDate:   startT,
		EndDate:     endT,
	}
	err = s.repo.FullUpdateSubscriptionsById(ctx, &sub)
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
func (s *SubscriptionService) GetSubscriptionsFilter(ctx context.Context, from, to, user_id, service_name string) (int, error) {
	fromT, err := date.ParseMonthDate(from)
	if err != nil {
		return 0, err
	}

	toT, err := date.ParseMonthDate(to)
	if err != nil {
		return 0, err
	}
	to = toT.Format("02-01-2006")
	from = fromT.Format("02-01-2006")
	total, err := s.repo.GetSubscriptionsFilter(ctx, from, to, user_id, service_name)
	if err != nil {
		return 0, fmt.Errorf("s.repo.GetSubscriptionsFilter: %w", err)
	}
	return total, nil
}
