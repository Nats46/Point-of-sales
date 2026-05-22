package form

type CreateSalesRequest struct {
	CustomerId    int64              `json:"customer_id"`
	CashierId     int64              `json:"cashier_id"`
	PaymentMethod string             `json:"payment_method" binding:"required"`
	Status        string             `json:"status" binding:"required"`
	SalesDetails  []CreateSalesDetail  `json:"sales_details" binding:"required,min=1"`
}

type CreateSalesDetail struct {
	ItemId   int64   `json:"item_id" binding:"required"`
	ItemName string  `json:"item_name" binding:"required"`
	ItemCode string  `json:"item_code" binding:"required"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Discount float64 `json:"discount"`
	Tax      float64 `json:"tax"`
}

type UpdateSalesRequest struct {
	CustomerId    int64              `json:"customer_id"`
	CashierId     int64              `json:"cashier_id"`
	PaymentMethod string             `json:"payment_method"`
	Status        string             `json:"status"`
	SalesDetails  []UpdateSalesDetail `json:"sales_details" binding:"required,min=1"`
}

type UpdateSalesDetail struct {
	SalesDetailId int64   `json:"sales_detail_id"`
	ItemId        int64   `json:"item_id"`
	ItemName      string  `json:"item_name"`
	ItemCode      string  `json:"item_code"`
	Quantity      float64 `json:"quantity"`
	Price         float64 `json:"price"`
	Discount      float64 `json:"discount"`
	Tax           float64 `json:"tax"`
}