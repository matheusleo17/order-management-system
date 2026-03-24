package events

type OrderCreatedEvent struct {
	OrderID string  `json:"order_id"`
	Total   float64 `json:"total"`
}
