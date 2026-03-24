package usecase

import (
	"context"

	"order-management-system/internal/repository"
)

type DeleteOrderUseCase struct {
	repo repository.OrderRepository
}

func NewDeleteOrderUseCase(repo repository.OrderRepository) *DeleteOrderUseCase {
	return &DeleteOrderUseCase{
		repo: repo,
	}
}

func (uc *DeleteOrderUseCase) Execute(ctx context.Context, id string) error {
	return uc.repo.DeleteOrder(ctx, id)
}
