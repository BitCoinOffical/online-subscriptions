package dto

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionDTO struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       float64   `json:"price" binding:"required"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
}
