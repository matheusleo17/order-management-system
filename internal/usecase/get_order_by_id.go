package usecase

import (
	"context"

	"order-management-system/internal/domain"
	"order-management-system/internal/repository"
)

type GetOrderByIdUseCase struct {
	repo repository.OrderRepository
}

func NewGetOrderByIdUseCase(repo repository.OrderRepository) *GetOrderByIdUseCase {
	return &GetOrderByIdUseCase{
		repo: repo,
	}
}

func (uc *GetOrderByIdUseCase) Execute(ctx context.Context, id string) (domain.Order, error) {
	return uc.repo.GetOrdersById(ctx, id)
}
