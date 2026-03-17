package repository

import (
	"context"
	"time"

	"order-management-system/internal/database"
	"order-management-system/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
)

type OrderRepository struct {
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

const collectionName = "orders"

func (r *OrderRepository) InsertOrder(ctx context.Context, order domain.Order) error {

	collection := database.DB.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, order)

	return err
}

func (r *OrderRepository) GetOrders(ctx context.Context) ([]domain.Order, error) {

	collection := database.DB.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	var orders []domain.Order

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) GetOrdersById(ctx context.Context, id string) (domain.Order, error) {

	collection := database.DB.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var order domain.Order

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &order); err != nil {
		return nil, err
	}

	return order, nil
}
