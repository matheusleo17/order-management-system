package usecase_test

import (
	"context"
	"errors"
	"testing"

	"order-management-system/internal/domain"
	"order-management-system/internal/usecase"
)

type mockOrderDeleteRepository struct {
	deleteErr error
}

func (m *mockOrderDeleteRepository) DeleteOrder(_ context.Context, _ string) error {
	return m.deleteErr
}

func (m *mockOrderDeleteRepository) InsertOrder(_ context.Context, _ domain.Order) error {
	return nil
}

func (m *mockOrderDeleteRepository) GetOrders(_ context.Context) ([]domain.Order, error) {
	return nil, nil
}

func (m *mockOrderDeleteRepository) GetOrdersById(_ context.Context, _ string) (domain.Order, error) {
	return domain.Order{}, nil
}

func (m *mockOrderDeleteRepository) UpdateOrder(_ context.Context, _ string, _ domain.Order) error {
	return nil
}

func TestDeleteOrder_Success(t *testing.T) {
	repo := &mockOrderDeleteRepository{}
	uc := usecase.NewDeleteOrderUseCase(repo)

	if err := uc.Execute(context.Background(), "order123"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDeleteOrder_NotFound(t *testing.T) {
	repo := &mockOrderDeleteRepository{deleteErr: errors.New("not found")}
	uc := usecase.NewDeleteOrderUseCase(repo)

	err := uc.Execute(context.Background(), "order123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
