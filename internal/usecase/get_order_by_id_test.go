package usecase_test

import (
	"context"
	"errors"
	"testing"

	"order-management-system/internal/domain"
	"order-management-system/internal/usecase"
)

type mockOrderGetByIdRepository struct {
	getErr error
}

func (m *mockOrderGetByIdRepository) DeleteOrder(_ context.Context, _ string) error {
	return m.getErr
}

func (m *mockOrderGetByIdRepository) InsertOrder(_ context.Context, _ domain.Order) error {
	return nil
}

func (m *mockOrderGetByIdRepository) GetOrders(_ context.Context) ([]domain.Order, error) {
	return nil, nil
}

func (m *mockOrderGetByIdRepository) GetOrdersById(_ context.Context, id string) (domain.Order, error) {
	return domain.Order{Id: id}, m.getErr
}

func (m *mockOrderGetByIdRepository) UpdateOrder(_ context.Context, _ string, _ domain.Order) error {
	return nil
}

func TestGetByIdOrders(t *testing.T) {
	repo := &mockOrderGetByIdRepository{}
	uc := usecase.NewGetOrderByIdUseCase(repo)

	order, err := uc.Execute(context.Background(), "order123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if order.Id != "order123" {
		t.Errorf("expected order id 'order123', got '%s'", order.Id)
	}
}

func TestGetByIdOrders_Error(t *testing.T) {
	repo := &mockOrderGetByIdRepository{getErr: errors.New("db error")}
	uc := usecase.NewGetOrderByIdUseCase(repo)

	_, err := uc.Execute(context.Background(), "order123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
