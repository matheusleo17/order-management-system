package service

import (
	"time"

	"github.com/google/uuid"

	"order-management-system/internal/domain"
)

var orders []domain.Order

func CreateOrder(order domain.Order) domain.Order {

	order.Id = uuid.New().String()
	order.Status = "created"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	var subtotal float64

	for _, item := range order.Items {
		subtotal += float64(item.Quantity) * item.Price
	}

	order.Subtotal = subtotal

	var discountTotal float64
	for _, d := range order.Discounts {
		discountTotal += d.Amount
	}

	var taxTotal float64
	for _, t := range order.Taxes {
		taxTotal += t.Amount
	}

	order.Total = subtotal - discountTotal + taxTotal

	orders = append(orders, order)

	return order
}

func GetOrders() []domain.Order {
	return orders
}
