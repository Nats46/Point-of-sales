package form

import "time"

type CreateBatchRequest struct {
	ItemCode    string    `json:"item_code" binding:"required"`
	ItemName    string    `json:"item_name" binding:"required"`
	BatchNumber string    `json:"batch_number" binding:"required"`
	Stock       float64   `json:"stock" binding:"required,gte=0"`
	Status      string    `json:"status" binding:"required"`
	ExpiredDate time.Time `json:"expired_date" binding:"required"`
	StockDate   time.Time `json:"stock_date" binding:"required"`
}

type UpdateBatchRequest struct {
	ItemCode    string    `json:"item_code"`
	ItemName    string    `json:"item_name"`
	BatchNumber string    `json:"batch_number"`
	Stock       float64   `json:"stock"`
	Status      string    `json:"status"`
	ExpiredDate time.Time `json:"expired_date"`
	StockDate   time.Time `json:"stock_date"`
}