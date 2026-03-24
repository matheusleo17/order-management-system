package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"order-management-system/internal/domain"
	"order-management-system/internal/repository"
)

type EventPublisher interface {
	Publish(eventName string, payload interface{}) error
}

type CreateOrderUseCase struct {
	repo      repository.OrderRepository
	publisher EventPublisher
}

func NewCreateOrderUseCase(
	repo repository.OrderRepository,
	publisher EventPublisher,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		repo:      repo,
		publisher: publisher,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, order *domain.Order) error {
	if len(order.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}

	order.Id = uuid.New().String()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	if order.Status == "" {
		order.Status = "pending"
	}

	if err := uc.repo.InsertOrder(ctx, *order); err != nil {
		return fmt.Errorf("failed to persist order: %w", err)
	}

	if err := uc.publisher.Publish("order.created", order); err != nil {
		return fmt.Errorf("order saved but failed to publish event: %w", err)
	}

	return nil
}
