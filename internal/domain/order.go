package domain

import "time"

type Order struct {
	Id        string      `json:"id"`
	Customer  Customer    `json:"customer"`
	Items     []OrderItem `json:"items"`
	Discounts []Discount  `json:"discounts"`
	Taxes     []Tax       `json:"taxes"`
	Payment   Payment     `json:"payment"`
	Shipment  Shipment    `json:"shipment"`
	Status    string      `json:"status"`
	Subtotal  float64     `json:"subtotal"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
