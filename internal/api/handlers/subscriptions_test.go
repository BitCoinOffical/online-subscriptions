package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BitCoinOffical/online-subscriptions/internal/api/handlers"
	mocks "github.com/BitCoinOffical/online-subscriptions/internal/api/handlers/mock"
	"github.com/BitCoinOffical/online-subscriptions/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/internal/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCreateSubscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockSubscriptionService(ctrl)

	handler := handlers.NewSubscriptionHandler(mockSvc, zap.NewNop())

	router := gin.New()
	router.POST("/api/v1/subscriptions", handler.CreateSubscription)

	testCases := []struct {
		name       string
		data       map[string]any
		wantStatus int
	}{
		{
			name: "success",
			data: map[string]any{
				"service_name": "Netflix",
				"price":        100,
				"user_id":      "550e8400-e29b-41d4-a716-446655440000",
				"start_date":   "2024-01-01",
				"end_date":     "2024-12-31",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing service_name",
			data: map[string]any{
				"price":      100,
				"user_id":    "550e8400-e29b-41d4-a716-446655440000",
				"start_date": "2024-01-01",
				"end_date":   "2024-12-31",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative price",
			data: map[string]any{
				"service_name": "Netflix",
				"price":        -1,
				"user_id":      "550e8400-e29b-41d4-a716-446655440000",
				"start_date":   "2024-01-01",
				"end_date":     "2024-12-31",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid uuid",
			data: map[string]any{
				"service_name": "Netflix",
				"price":        100,
				"user_id":      "not-a-uuid",
				"start_date":   "2024-01-01",
				"end_date":     "2024-12-31",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			data:       map[string]any{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantStatus == http.StatusCreated {
				mockSvc.EXPECT().CreateSubscription(gomock.Any(), gomock.Any()).Return(nil)
			}

			data, err := json.Marshal(tc.data)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/subscriptions", strings.NewReader(string(data)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

func TestUpdateSubscriptionsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockSubscriptionService(ctrl)

	handler := handlers.NewSubscriptionHandler(mockSvc, zap.NewNop())

	router := gin.New()
	router.PATCH("/api/v1/subscriptions/:id", handler.UpdateSubscriptionsById)
	testCases := []struct {
		name       string
		url        string
		data       map[string]any
		wantStatus int
	}{
		{
			name: "success",
			url:  "/api/v1/subscriptions/1",
			data: map[string]any{
				"service_name": "Netflix",
				"price":        200,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			url:        "/api/v1/subscriptions/abc",
			data:       map[string]any{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid body",
			url:        "/api/v1/subscriptions/1",
			data:       nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "internal server error",
			url:  "/api/v1/subscriptions/2",
			data: map[string]any{
				"service_name": "Netflix",
				"price":        200,
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "not found error",
			url:  "/api/v1/subscriptions/2444",
			data: map[string]any{
				"service_name": "Netflix",
				"price":        200,
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			switch tc.wantStatus {
			case http.StatusOK:
				mockSvc.EXPECT().UpdateSubscriptionsById(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			case http.StatusNotFound:
				mockSvc.EXPECT().UpdateSubscriptionsById(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.ErrNotFound)
			case http.StatusInternalServerError:
				mockSvc.EXPECT().UpdateSubscriptionsById(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			}

			data, err := json.Marshal(tc.data)

			req := httptest.NewRequest(http.MethodPatch, tc.url, strings.NewReader(string(data)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

func TestFullUpdateSubscriptionsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockSubscriptionService(ctrl)

	handler := handlers.NewSubscriptionHandler(mockSvc, zap.NewNop())

	router := gin.New()
	router.PUT("/api/v1/subscriptions/:id", handler.FullUpdateSubscriptionsById)

	testCases := []struct {
		name       string
		url        string
		data       map[string]any
		wantStatus int
	}{
		{
			name: "success",
			url:  "/api/v1/subscriptions/1",
			data: map[string]any{
				"service_name": "Netflix",
				"price":        100,
				"user_id":      "550e8400-e29b-41d4-a716-446655440000",
				"start_date":   "2024-01-01",
				"end_date":     "2024-12-31",
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			url:        "/api/v1/subscriptions/abc",
			data:       map[string]any{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid body",
			url:        "/api/v1/subscriptions/1",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing required field",
			url:  "/api/v1/subscriptions/1",
			data: map[string]any{
				"price":      100,
				"user_id":    "550e8400-e29b-41d4-a716-446655440000",
				"start_date": "2024-01-01",
				"end_date":   "2024-12-31",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "not found",
			url:  "/api/v1/subscriptions/999",
			data: map[string]any{
				"service_name": "Netflix",
				"price":        100,
				"user_id":      "550e8400-e29b-41d4-a716-446655440000",
				"start_date":   "2024-01-01",
				"end_date":     "2024-12-31",
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal server error",
			url:  "/api/v1/subscriptions/2",
			data: map[string]any{
				"service_name": "Netflix",
				"price":        100,
				"user_id":      "550e8400-e29b-41d4-a716-446655440000",
				"start_date":   "2024-01-01",
				"end_date":     "2024-12-31",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.wantStatus {
			case http.StatusOK:
				mockSvc.EXPECT().FullUpdateSubscriptionsById(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			case http.StatusNotFound:
				mockSvc.EXPECT().FullUpdateSubscriptionsById(gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.ErrNotFound)
			case http.StatusInternalServerError:
				mockSvc.EXPECT().FullUpdateSubscriptionsById(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			}

			data, err := json.Marshal(tc.data)
			req := httptest.NewRequest(http.MethodPut, tc.url, strings.NewReader(string(data)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

func TestGetSubscriptionsById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockSubscriptionService(ctrl)

	handler := handlers.NewSubscriptionHandler(mockSvc, zap.NewNop())

	router := gin.New()
	router.GET("/api/v1/subscriptions/:id", handler.GetSubscriptionsById)

	testCases := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "success",
			url:        "/api/v1/subscriptions/1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			url:        "/api/v1/subscriptions/999",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id",
			url:        "/api/v1/subscriptions/abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "internal server error",
			url:        "/api/v1/subscriptions/2",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.wantStatus {
			case http.StatusOK:
				mockSvc.EXPECT().GetSubscriptionsById(gomock.Any(), 1).Return(&models.Subscription{
					ID: 1, ServiceName: "Netflix",
				}, nil)
			case http.StatusNotFound:
				mockSvc.EXPECT().GetSubscriptionsById(gomock.Any(), 999).Return(nil, domain.ErrNotFound)
			case http.StatusInternalServerError:
				mockSvc.EXPECT().GetSubscriptionsById(gomock.Any(), 2).Return(nil, errors.New("db error"))
			}

			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.NoError(t, err)
			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

func TestDeleteSubscriptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockSubscriptionService(ctrl)

	handler := handlers.NewSubscriptionHandler(mockSvc, zap.NewNop())

	router := gin.New()
	router.DELETE("/api/v1/subscriptions/:id", handler.DeleteSubscriptions)

	testCases := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "success",
			url:        "/api/v1/subscriptions/1",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid id",
			url:        "/api/v1/subscriptions/abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "internal server error",
			url:        "/api/v1/subscriptions/2",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "not found",
			url:        "/api/v1/subscriptions/999",
			wantStatus: http.StatusNotFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.wantStatus {
			case http.StatusNoContent:
				mockSvc.EXPECT().DeleteSubscriptions(gomock.Any(), 1).Return(nil)
			case http.StatusNotFound:
				mockSvc.EXPECT().DeleteSubscriptions(gomock.Any(), 999).Return(domain.ErrNotFound)
			case http.StatusInternalServerError:
				mockSvc.EXPECT().DeleteSubscriptions(gomock.Any(), 2).Return(errors.New("db error"))
			}

			req, err := http.NewRequest(http.MethodDelete, tc.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.NoError(t, err)
			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

func TestGetSubscriptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockSubscriptionService(ctrl)

	handler := handlers.NewSubscriptionHandler(mockSvc, zap.NewNop())

	router := gin.New()
	router.GET("/api/v1/subscriptions/", handler.GetSubscriptions)

	testCases := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "success",
			url:        "/api/v1/subscriptions/",
			wantStatus: http.StatusOK,
		},
		{
			name:       "internal server error",
			url:        "/api/v1/subscriptions/",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.wantStatus {
			case http.StatusOK:
				mockSvc.EXPECT().GetSubscriptions(gomock.Any()).Return([]models.Subscription{}, nil)
			case http.StatusInternalServerError:
				mockSvc.EXPECT().GetSubscriptions(gomock.Any()).Return(nil, errors.New("db error"))
			}

			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.NoError(t, err)
			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

func TestGetSubscriptionsFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockSubscriptionService(ctrl)

	handler := handlers.NewSubscriptionHandler(mockSvc, zap.NewNop())

	router := gin.New()
	router.GET("/api/v1/subscriptions", handler.GetSubscriptionsFilter)

	testCases := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "success all filters",
			url:        "/api/v1/subscriptions?user_id=550e8400-e29b-41d4-a716-446655440000&service_name=Netflix&from=2024-01-01&to=2024-12-31",
			wantStatus: http.StatusOK,
		},
		{
			name:       "success no filters",
			url:        "/api/v1/subscriptions",
			wantStatus: http.StatusOK,
		},
		{
			name:       "success only user_id",
			url:        "/api/v1/subscriptions?user_id=550e8400-e29b-41d4-a716-446655440000",
			wantStatus: http.StatusOK,
		},
		{
			name:       "internal server error",
			url:        "/api/v1/subscriptions?service_name=error",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.wantStatus {
			case http.StatusOK:
				mockSvc.EXPECT().GetSubscriptionsFilter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(100, nil)
			case http.StatusInternalServerError:
				mockSvc.EXPECT().GetSubscriptionsFilter(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(0, errors.New("db error"))
			}

			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.NoError(t, err)
			assert.Equal(t, tc.wantStatus, w.Code)
		})
	}
}
