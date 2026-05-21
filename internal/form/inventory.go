package form

type CreateInventoryRequest struct {
	ItemCode string  `json:"item_code" binding:"required"`
	ItemName string  `json:"item_name" binding:"required"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Unit     string  `json:"unit" binding:"required"`
}

type UpdateInventoryRequest struct {
	ItemCode string  `json:"item_code"`
	ItemName string  `json:"item_name"`
	Price    float64 `json:"price"`
	Unit     string  `json:"unit"`
}