package model

import "time"

type Inventory struct {
	Id        int64     `db:"id" json:"id"`
	ItemCode  string    `db:"item_code" json:"item_code"`
	ItemName  string    `db:"item_name" json:"item_name"`
	Price     float64   `db:"price" json:"price"`
	Unit      string    `db:"unit" json:"unit"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy string    `db:"created_by" json:"created_by"`
	UpdatedBy string    `db:"updated_by" json:"updated_by"`
}

type Batch struct {
	Id          int64     `db:"id" json:"id"`
	ItemCode    string    `db:"item_code" json:"item_code"`
	ItemName    string    `db:"item_name" json:"item_name"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
	Stock       float64   `db:"stock" json:"stock"`
	Status      string    `db:"status" json:"status"`
	BatchNumber string    `db:"batch_number" json:"batch_number"`
	ExpiredDate time.Time `db:"expired_date" json:"expired_date"`
	StockDate   time.Time `db:"stock_date" json:"stock_date"`
}