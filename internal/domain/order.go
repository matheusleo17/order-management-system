package domain

import "time"

type Order struct {
	Id        string      `json:"id" bson:"_id"`
	Customer  Customer    `json:"customer" bson:"customer"`
	Items     []OrderItem `json:"items" bson:"items"`
	Discounts []Discount  `json:"discounts" bson:"discounts"`
	Taxes     []Tax       `json:"taxes" bson:"taxes"`
	Payment   Payment     `json:"payment" bson:"payment"`
	Shipment  Shipment    `json:"shipment" bson:"shipment"`
	Status    string      `json:"status" bson:"status"`
	Subtotal  float64     `json:"subtotal" bson:"subtotal"`
	Total     float64     `json:"total" bson:"total"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" bson:"updated_at"`
}
