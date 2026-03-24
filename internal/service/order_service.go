package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"order-management-system/internal/domain"
	"order-management-system/internal/events"
	"order-management-system/internal/messaging"
	"order-management-system/internal/repository"
)

type OrderService struct {
	repo      *repository.OrderRepository
	publisher *messaging.RabbitMqPublisher
}

func NewOrderService(repo *repository.OrderRepository, publisher *messaging.RabbitMqPublisher) *OrderService {
	return &OrderService{
		repo:      repo,
		publisher: publisher,
	}
}
func (s *OrderService) CreateOrder(ctx context.Context, order domain.Order) (domain.Order, error) {

	order.Id = uuid.New().String()
	order.Status = "created"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	event := events.OrderCreatedEvent{
		OrderID: order.Id,
		Total:   order.Total,
	}

	_ = s.publisher.Publish("order.created", event)
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

	err := s.repo.InsertOrder(ctx, order)

	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
func (s *OrderService) GetOrders(ctx context.Context) ([]domain.Order, error) {
	return s.repo.GetOrders(ctx)
}

func (s *OrderService) GetOrderById(ctx context.Context, id string) (domain.Order, error) {
	return s.repo.GetOrdersById(ctx, id)
}

func (s *OrderService) UpdateOrder(ctx context.Context, id string, order domain.Order) error {
	return s.repo.UpdateOrder(ctx, id, order)
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	return s.repo.DeleteOrder(ctx, id)
}
