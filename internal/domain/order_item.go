package domain

type OrderItem struct {
	ProductId string  `json:"product_id"`
	Product   string  `json:"product"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
