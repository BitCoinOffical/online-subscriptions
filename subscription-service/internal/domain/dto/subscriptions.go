package dto

import (
	"github.com/google/uuid"
)

// SubscriptionDTO represents subscription creation/update request
type SubscriptionDTO struct {
	ServiceName string    `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price       int       `json:"price" binding:"required,gte=0" example:"400"`
	UserID      uuid.UUID `json:"user_id" binding:"required" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date" binding:"required" example:"07-2025"`
	EndDate     string    `json:"end_date" binding:"omitempty" example:"12-2025"`
}

// PatchSubscriptionDTO represents partial subscription update request
type PatchSubscriptionDTO struct {
	ServiceName string    `json:"service_name" binding:"omitempty" example:"Netflix"`
	Price       *int      `json:"price" binding:"omitempty,gte=0" example:"500"`
	UserID      uuid.UUID `json:"user_id" binding:"omitempty" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date" binding:"omitempty" example:"01-2026"`
	EndDate     string    `json:"end_date" binding:"omitempty" example:"06-2026"`
}

func (d PatchSubscriptionDTO) IsEmpty() bool {
	return d.ServiceName == "" &&
		d.Price == nil &&
		d.UserID == uuid.Nil &&
		d.StartDate == "" &&
		d.EndDate == ""
}
