package usecase

import (
	"context"
	"fmt"

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

	// 2. Persistir
	err := uc.repo.InsertOrder(ctx, *order)
	if err != nil {
		return err
	}

	// 3. Publicar evento
	err = uc.publisher.Publish("order.created", order)
	if err != nil {
		return err
	}

	return nil
}
