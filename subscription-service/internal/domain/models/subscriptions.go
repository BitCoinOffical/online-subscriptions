package models

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a user subscription
type Subscription struct {
	ID          int        `json:"id" example:"1"`
	ServiceName string     `json:"service_name" example:"Yandex Plus"`
	Price       int        `json:"price" example:"400"`
	UserID      uuid.UUID  `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time  `json:"start_date" example:"2025-07-01T00:00:00Z"`
	EndDate     *time.Time `json:"end_date,omitempty" example:"2025-12-31T00:00:00Z"`
	Created_at  time.Time  `json:"created_at" example:"2026-05-12T10:00:00Z"`
	Updated_at  time.Time  `json:"updated_at" example:"2026-05-12T10:00:00Z"`
}
type PatchSubscription struct {
	ID          int
	ServiceName string
	Price       *int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
	Created_at  time.Time
	Updated_at  time.Time
}
