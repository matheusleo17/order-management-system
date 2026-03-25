package usecase_test

import (
	"context"
	"errors"
	"testing"

	"order-management-system/internal/domain"
	"order-management-system/internal/usecase"
)

type mockOrderGetRepository struct {
	getErr error
}

func (m *mockOrderGetRepository) GetOrdersById(_ context.Context, _ string) (domain.Order, error) {
	return domain.Order{}, m.getErr
}

func (m *mockOrderGetRepository) InsertOrder(_ context.Context, _ domain.Order) error {
	return nil
}

func (m *mockOrderGetRepository) GetOrders(_ context.Context) ([]domain.Order, error) {
	return []domain.Order{{Id: "order123"}}, m.getErr
}

func (m *mockOrderGetRepository) UpdateOrder(_ context.Context, _ string, _ domain.Order) error {
	return nil
}
func (m *mockOrderGetRepository) DeleteOrder(_ context.Context, _ string) error {
	return nil
}

func TestGetAllOrders(t *testing.T) {
	repo := &mockOrderGetRepository{}
	uc := usecase.NewGetOrdersUseCase(repo)

	orders, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(orders) != 1 {
		t.Errorf("expected 1 order, got %d", len(orders))
	}
}

func TestGetAllOrders_Error(t *testing.T) {
	repo := &mockOrderGetRepository{getErr: errors.New("db error")}
	uc := usecase.NewGetOrdersUseCase(repo)

	_, err := uc.Execute(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
