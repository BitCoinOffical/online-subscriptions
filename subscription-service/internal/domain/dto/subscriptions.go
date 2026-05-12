package dto

import (
	"github.com/google/uuid"
)

// SubscriptionDTO represents subscription creation/update request
type SubscriptionDTO struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,gte=0"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     string    `json:"end_date" binding:"omitempty"`
}

// PatchSubscriptionDTO represents partial subscription update request
type PatchSubscriptionDTO struct {
	ServiceName string    `json:"service_name" binding:"omitempty"`
	Price       *int      `json:"price" binding:"omitempty,gte=0"`
	UserID      uuid.UUID `json:"user_id" binding:"omitempty"`
	StartDate   string    `json:"start_date" binding:"omitempty"`
	EndDate     string    `json:"end_date" binding:"omitempty"`
}

func (d PatchSubscriptionDTO) IsEmpty() bool {
	return d.ServiceName == "" &&
		d.Price == nil &&
		d.UserID == uuid.Nil &&
		d.StartDate == "" &&
		d.EndDate == ""
}
