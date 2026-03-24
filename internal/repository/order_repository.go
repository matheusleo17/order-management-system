package repository

import (
	"context"

	"order-management-system/internal/domain"
)

type OrderRepository interface {
	InsertOrder(ctx context.Context, order domain.Order) error
	GetOrders(ctx context.Context) ([]domain.Order, error)
	GetOrdersById(ctx context.Context, id string) (domain.Order, error)
	UpdateOrder(ctx context.Context, id string, order domain.Order) error
	DeleteOrder(ctx context.Context, id string) error
}
