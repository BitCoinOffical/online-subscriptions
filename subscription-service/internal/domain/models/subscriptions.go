package models

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a user subscription
type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
	Created_at  time.Time
	Updated_at  time.Time
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
