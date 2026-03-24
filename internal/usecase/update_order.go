package usecase

import (
	"context"

	"order-management-system/internal/domain"
	"order-management-system/internal/repository"
)

type UpdateOrderUseCase struct {
	repo repository.OrderRepository
}

func NewUpdateOrderUseCase(repo repository.OrderRepository) *UpdateOrderUseCase {
	return &UpdateOrderUseCase{
		repo: repo,
	}
}

func (uc *UpdateOrderUseCase) Execute(ctx context.Context, id string, order *domain.Order) error {
	return uc.repo.UpdateOrder(ctx, id, *order)
}
