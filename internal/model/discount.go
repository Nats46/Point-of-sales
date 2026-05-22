package model

import "time"

type Discount struct {
	Id           int64
	DiscountCode string
	Value        float64
	DiscountType string    `db:"discount_type" json:"discount_type"`
	StartDate    time.Time `db:"start_date" json:"start_date"`
	EndDate      time.Time `db:"end_date" json:"end_date"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy    string    `db:"created_by" json:"created_by"`
	UpdatedBy    string    `db:"updated_by" json:"updated_by"`
}
