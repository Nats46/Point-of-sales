package model

import "time"

type SalesHeader struct {
	TransactionId int64         `db:"transaction_id" json:"transaction_id"`
	TransactionNo string        `db:"transaction_no" json:"transaction_no"`
	CustomerId    int64         `db:"customer_id" json:"customer_id"`
	CashierId     int64         `db:"cashier_id" json:"cashier_id"`
	Subtotal      float64       `db:"subtotal" json:"subtotal"`
	Discount      float64       `db:"discount" json:"discount"`
	Tax           float64       `db:"tax" json:"tax"`
	GrandTotal    float64       `db:"grand_total" json:"grand_total"`
	PaymentMethod string        `db:"payment_method" json:"payment_method"`
	Status        string        `db:"status" json:"status"`
	CreatedAt     time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" json:"updated_at"`
	CreatedBy     string        `db:"created_by" json:"created_by"`
	UpdatedBy     string        `db:"updated_by" json:"updated_by"`
	SalesDetails  []*SalesDetail `json:"sales_details"`
}

type SalesDetail struct {
	SalesDetailId int64     `db:"sales_detail_id" json:"sales_detail_id"`
	TransactionId int64     `db:"transaction_id" json:"transaction_id"`
	ItemId        int64     `db:"item_id" json:"item_id"`
	ItemName      string    `db:"item_name" json:"item_name"`
	ItemCode      string    `db:"item_code" json:"item_code"`
	Quantity      float64   `db:"quantity" json:"quantity"`
	Price         float64   `db:"price" json:"price"`
	Subtotal      float64   `db:"subtotal" json:"subtotal"`
	Discount      float64   `db:"discount" json:"discount"`
	Tax           float64   `db:"tax" json:"tax"`
	GrandTotal    float64   `db:"grand_total" json:"grand_total"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	CreatedBy     string    `db:"created_by" json:"created_by"`
	UpdatedBy     string    `db:"updated_by" json:"updated_by"`
}
