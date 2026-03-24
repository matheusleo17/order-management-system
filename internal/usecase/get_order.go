package usecase

import (
	"context"

	"order-management-system/internal/domain"
	"order-management-system/internal/repository"
)

type GetOrdersUseCase struct {
	repo repository.OrderRepository
}

func NewGetOrdersUseCase(repo repository.OrderRepository) *GetOrdersUseCase {
	return &GetOrdersUseCase{
		repo: repo,
	}
}

func (uc *GetOrdersUseCase) Execute(ctx context.Context) ([]domain.Order, error) {
	return uc.repo.GetOrders(ctx)
}
