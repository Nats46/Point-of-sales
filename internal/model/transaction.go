package model

import "time"

type Transaction struct {
	ID          int64     `db:"id" json:"id"`
	Transaction string    `db:"transaction" json:"transaction"`
	Requester   int64     `db:"requester" json:"requester"`
	Approver    int64     `db:"approver" json:"approver"`
	Status      string    `db:"status" json:"status"`
	GrandTotal  float64   `db:"grand_total" json:"grand_total"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}

type TransactionItems struct {
	ID          int64     `db:"id" json:"id"`
	Transaction string    `db:"transaction" json:"transaction"`
	InventoryId string    `db:"item" json:"item"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Price       float64   `db:"price" json:"price"`
	Total       float64   `db:"total" json:"total"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	UpdatedBy   string    `db:"updated_by" json:"updated_by"`
}
