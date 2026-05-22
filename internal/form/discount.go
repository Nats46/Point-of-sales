package form

import "time"

type CreateDiscountRequest struct {
	DiscountCode string    `json:"discount_code" binding:"required"`
	Value        float64   `json:"value" binding:"required"`
	DiscountType string    `json:"discount_type" binding:"required"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required"`
}

type UpdateDiscountRequest struct {
	DiscountCode string    `json:"discount_code"`
	Value        float64   `json:"value"`
	DiscountType string    `json:"discount_type"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	IsActive     bool      `json:"is_active"`
}