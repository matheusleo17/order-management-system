package usecase_test

import (
	"context"
	"errors"
	"testing"

	"order-management-system/internal/domain"
	"order-management-system/internal/usecase"
)

// --- Mocks ---

type mockOrderRepository struct {
	insertErr error
	deleteErr error
	inserted  *domain.Order
}

func (m *mockOrderRepository) InsertOrder(_ context.Context, order domain.Order) error {
	m.inserted = &order
	return m.insertErr
}

func (m *mockOrderRepository) GetOrders(_ context.Context) ([]domain.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) GetOrdersById(_ context.Context, _ string) (domain.Order, error) {
	return domain.Order{}, nil
}

func (m *mockOrderRepository) UpdateOrder(_ context.Context, _ string, _ domain.Order) error {
	return nil
}

func (m *mockOrderRepository) DeleteOrder(_ context.Context, _ string) error {
	return nil
}

type mockPublisher struct {
	publishErr error
	published  bool
}

func (m *mockPublisher) Publish(_ string, _ interface{}) error {
	m.published = true
	return m.publishErr
}

// --- Tests ---

func TestCreateOrder_Success(t *testing.T) {
	repo := &mockOrderRepository{}
	pub := &mockPublisher{}

	uc := usecase.NewCreateOrderUseCase(repo, pub)

	order := &domain.Order{
		Items: []domain.OrderItem{
			{ProductId: "p1", Product: "Widget", Quantity: 2, Price: 10.0},
		},
	}

	if err := uc.Execute(context.Background(), order); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if order.Id == "" {
		t.Error("expected order ID to be generated")
	}

	if order.Status != "pending" {
		t.Errorf("expected status 'pending', got '%s'", order.Status)
	}

	if order.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if !pub.published {
		t.Error("expected event to be published")
	}
}

func TestCreateOrder_NoItems(t *testing.T) {
	repo := &mockOrderRepository{}
	pub := &mockPublisher{}

	uc := usecase.NewCreateOrderUseCase(repo, pub)

	err := uc.Execute(context.Background(), &domain.Order{})
	if err == nil {
		t.Fatal("expected error for empty items, got nil")
	}
}

func TestCreateOrder_RepositoryError(t *testing.T) {
	repo := &mockOrderRepository{insertErr: errors.New("db error")}
	pub := &mockPublisher{}

	uc := usecase.NewCreateOrderUseCase(repo, pub)

	order := &domain.Order{
		Items: []domain.OrderItem{{ProductId: "p1", Quantity: 1, Price: 5.0}},
	}

	if err := uc.Execute(context.Background(), order); err == nil {
		t.Fatal("expected error from repository, got nil")
	}
}

func TestCreateOrder_PublisherError(t *testing.T) {
	repo := &mockOrderRepository{}
	pub := &mockPublisher{publishErr: errors.New("rabbit down")}

	uc := usecase.NewCreateOrderUseCase(repo, pub)

	order := &domain.Order{
		Items: []domain.OrderItem{{ProductId: "p1", Quantity: 1, Price: 5.0}},
	}

	if err := uc.Execute(context.Background(), order); err == nil {
		t.Fatal("expected error from publisher, got nil")
	}
}

func TestDeleteOrder_Success(t *testing.T) {
	repo := &mockOrderRepository{}
	uc := usecase.NewDeleteOrderUseCase(repo)

	if err := uc.Execute(context.Background(), "order123"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestDeleteOrder_NotFound(t *testing.T) {
	repo := &mockOrderRepository{deleteErr: errors.New("not found")}
	uc := usecase.NewDeleteOrderUseCase(repo)

	err := uc.Execute(context.Background(), "order123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
