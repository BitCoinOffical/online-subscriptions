package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain/dto"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain/models"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/interfaces/http/services"
	mocks_services "github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/interfaces/http/services/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_services.NewMockSubscriptionRepo(ctrl)

	services := services.NewSubscriptionService(mockRepo)

	testCases := []struct {
		name    string
		dto     *dto.SubscriptionDTO
		setup   func(mockRepo *mocks_services.MockSubscriptionRepo)
		wantErr bool
	}{
		{
			name: "success",
			dto: &dto.SubscriptionDTO{
				ServiceName: "Netflix",
				Price:       100,
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				StartDate:   "01-2024",
				EndDate:     "12-2024",
			},
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().CreateSubscription(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid start_date",
			dto: &dto.SubscriptionDTO{
				ServiceName: "Netflix",
				Price:       100,
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				StartDate:   "invalid-date",
				EndDate:     "12-2024",
			},
			setup:   func(mockRepo *mocks_services.MockSubscriptionRepo) {},
			wantErr: true,
		},
		{
			name: "invalid end_date",
			dto: &dto.SubscriptionDTO{
				ServiceName: "Netflix",
				Price:       100,
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				StartDate:   "01-2024",
				EndDate:     "invalid-date",
			},
			setup:   func(mockRepo *mocks_services.MockSubscriptionRepo) {},
			wantErr: true,
		},
		{
			name: "repo error",
			dto: &dto.SubscriptionDTO{
				ServiceName: "Netflix",
				Price:       100,
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				StartDate:   "01-2024",
				EndDate:     "12-2024",
			},
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().CreateSubscription(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(mockRepo)

			err := services.CreateSubscription(context.Background(), tc.dto)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateSubscriptionsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_services.NewMockSubscriptionRepo(ctrl)

	services := services.NewSubscriptionService(mockRepo)

	testCases := []struct {
		name    string
		dto     *dto.PatchSubscriptionDTO
		id      int
		setup   func(mockRepo *mocks_services.MockSubscriptionRepo)
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			dto: &dto.PatchSubscriptionDTO{
				StartDate: "01-2024",
				EndDate:   "12-2024",
			},
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().UpdateSubscriptionsById(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid start_date",
			id:   1,
			dto: &dto.PatchSubscriptionDTO{
				StartDate: "invalid-date",
				EndDate:   "12-2024",
			},
			setup:   func(mockRepo *mocks_services.MockSubscriptionRepo) {},
			wantErr: true,
		},
		{
			name: "invalid end_date",
			id:   1,
			dto: &dto.PatchSubscriptionDTO{
				StartDate: "01-2024",
				EndDate:   "invalid-date",
			},
			setup:   func(mockRepo *mocks_services.MockSubscriptionRepo) {},
			wantErr: true,
		},
		{
			name: "repo error",
			id:   1,
			dto: &dto.PatchSubscriptionDTO{
				StartDate: "01-2024",
				EndDate:   "12-2024",
			},
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().UpdateSubscriptionsById(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "not found",
			id:   999,
			dto: &dto.PatchSubscriptionDTO{
				StartDate: "01-2024",
				EndDate:   "12-2024",
			},
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().UpdateSubscriptionsById(gomock.Any(), gomock.Any()).Return(domain.ErrNotFound)
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(mockRepo)

			err := services.UpdateSubscriptionsById(context.Background(), tc.dto, tc.id)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFullUpdateSubscriptionsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_services.NewMockSubscriptionRepo(ctrl)

	services := services.NewSubscriptionService(mockRepo)

	testCases := []struct {
		name    string
		dto     *dto.SubscriptionDTO
		id      int
		setup   func(mockRepo *mocks_services.MockSubscriptionRepo)
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			dto: &dto.SubscriptionDTO{
				ServiceName: "Netflix",
				Price:       100,
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				StartDate:   "01-2024",
				EndDate:     "12-2024",
			},
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().FullUpdateSubscriptionsById(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid start_date",
			id:   1,
			dto: &dto.SubscriptionDTO{
				ServiceName: "Netflix",
				Price:       100,
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				StartDate:   "invalid-date",
				EndDate:     "12-2024",
			},
			setup:   func(mockRepo *mocks_services.MockSubscriptionRepo) {},
			wantErr: true,
		},
		{
			name: "invalid end_date",
			id:   1,
			dto: &dto.SubscriptionDTO{
				ServiceName: "Netflix",
				Price:       100,
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				StartDate:   "01-2024",
				EndDate:     "invalid-date",
			},
			setup:   func(mockRepo *mocks_services.MockSubscriptionRepo) {},
			wantErr: true,
		},
		{
			name: "not found",
			id:   999,
			dto: &dto.SubscriptionDTO{
				ServiceName: "Netflix",
				Price:       100,
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				StartDate:   "01-2024",
				EndDate:     "12-2024",
			},
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().FullUpdateSubscriptionsById(gomock.Any(), gomock.Any()).Return(domain.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "repo error",
			id:   1,
			dto: &dto.SubscriptionDTO{
				ServiceName: "Netflix",
				Price:       100,
				UserID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				StartDate:   "01-2024",
				EndDate:     "12-2024",
			},
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().FullUpdateSubscriptionsById(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(mockRepo)

			err := services.FullUpdateSubscriptionsById(context.Background(), tc.dto, tc.id)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetSubscriptionsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_services.NewMockSubscriptionRepo(ctrl)

	services := services.NewSubscriptionService(mockRepo)

	testCases := []struct {
		name    string
		id      int
		setup   func(mockRepo *mocks_services.MockSubscriptionRepo)
		want    *models.Subscription
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().GetSubscriptionsById(gomock.Any(), 1).Return(&models.Subscription{
					ID:          1,
					ServiceName: "Netflix",
				}, nil)
			},
			want:    &models.Subscription{ID: 1, ServiceName: "Netflix"},
			wantErr: false,
		},
		{
			name: "repo error",
			id:   99,
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().GetSubscriptionsById(gomock.Any(), 99).Return(nil, domain.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(mockRepo)

			models, err := services.GetSubscriptionsById(context.Background(), tc.id)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, models)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, models)
			}
		})
	}
}

func TestDeleteSubscriptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_services.NewMockSubscriptionRepo(ctrl)

	services := services.NewSubscriptionService(mockRepo)

	testCases := []struct {
		name    string
		id      int
		setup   func(mockRepo *mocks_services.MockSubscriptionRepo)
		wantErr bool
	}{
		{
			name: "success",
			id:   1,
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().DeleteSubscriptions(gomock.Any(), 1).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repo error",
			id:   99,
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().DeleteSubscriptions(gomock.Any(), 99).Return(domain.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(mockRepo)

			err := services.DeleteSubscriptions(context.Background(), tc.id)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetSubscriptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_services.NewMockSubscriptionRepo(ctrl)
	service := services.NewSubscriptionService(mockRepo)

	testCases := []struct {
		name    string
		setup   func(mockRepo *mocks_services.MockSubscriptionRepo)
		want    []models.Subscription
		wantErr bool
	}{
		{
			name: "success",
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().
					GetSubscriptions(gomock.Any(), 10, 0).
					Return([]models.Subscription{
						{
							ID:          1,
							ServiceName: "Netflix",
						},
					}, nil)
			},
			want: []models.Subscription{
				{ID: 1, ServiceName: "Netflix"},
			},
			wantErr: false,
		},
		{
			name: "repo error",
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().
					GetSubscriptions(gomock.Any(), 10, 0).
					Return(nil, domain.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(mockRepo)

			got, err := service.GetSubscriptions(context.Background(), 10, 0)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestGetSubscriptionsFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_services.NewMockSubscriptionRepo(ctrl)

	services := services.NewSubscriptionService(mockRepo)

	testCases := []struct {
		name        string
		from        string
		to          string
		userID      string
		serviceName string
		setup       func(mockRepo *mocks_services.MockSubscriptionRepo)
		want        int
		wantErr     bool
	}{
		{
			name:        "success",
			from:        "01-2024",
			to:          "12-2024",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			serviceName: "Netflix",
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().GetSubscriptionsFilter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(100, nil)
			},
			want:    100,
			wantErr: false,
		},
		{
			name:        "invalid from",
			from:        "invalid-date",
			to:          "12-2024",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			serviceName: "Netflix",
			setup:       func(mockRepo *mocks_services.MockSubscriptionRepo) {},
			want:        0,
			wantErr:     true,
		},
		{
			name:        "invalid to",
			from:        "01-2024",
			to:          "invalid-date",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			serviceName: "Netflix",
			setup:       func(mockRepo *mocks_services.MockSubscriptionRepo) {},
			want:        0,
			wantErr:     true,
		},
		{
			name:        "empty filters",
			from:        "",
			to:          "",
			userID:      "",
			serviceName: "",
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().GetSubscriptionsFilter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(0, nil)
			},
			want:    0,
			wantErr: false,
		},
		{
			name:        "repo error",
			from:        "01-2024",
			to:          "12-2024",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			serviceName: "Netflix",
			setup: func(mockRepo *mocks_services.MockSubscriptionRepo) {
				mockRepo.EXPECT().GetSubscriptionsFilter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(0, errors.New("db error"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(mockRepo)

			total, err := services.GetSubscriptionsFilter(context.Background(), tc.from, tc.to, tc.userID, tc.serviceName)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Zero(t, total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, total)
			}
		})
	}
}
