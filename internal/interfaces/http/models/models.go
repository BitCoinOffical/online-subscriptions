package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int
	ServiceName string
	Price       float64
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     time.Time
	Created_at  time.Time
	Updated_at  time.Time
}
