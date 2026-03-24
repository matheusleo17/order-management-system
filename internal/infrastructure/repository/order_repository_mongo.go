package repository

import (
	"context"
	"fmt"
	"time"

	"order-management-system/internal/domain"
	"order-management-system/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepositoryMongo struct{}

func NewOrderRepositoryMongo() *OrderRepositoryMongo {
	return &OrderRepositoryMongo{}
}

const collectionName = "orders"

func (r *OrderRepositoryMongo) InsertOrder(ctx context.Context, order domain.Order) error {
	collection := database.GetDB().Collection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	return nil
}

func (r *OrderRepositoryMongo) GetOrders(ctx context.Context) ([]domain.Order, error) {
	collection := database.GetDB().Collection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find orders: %w", err)
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, fmt.Errorf("failed to decode orders: %w", err)
	}

	return orders, nil
}

func (r *OrderRepositoryMongo) GetOrdersById(ctx context.Context, id string) (domain.Order, error) {
	collection := database.GetDB().Collection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var order domain.Order
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Order{}, fmt.Errorf("order not found: %s", id)
		}
		return domain.Order{}, fmt.Errorf("failed to find order: %w", err)
	}

	return order, nil
}

type OrderUpdateFields struct {
	Items     []domain.OrderItem `bson:"items"`
	Discounts []domain.Discount  `bson:"discounts"`
	Taxes     []domain.Tax       `bson:"taxes"`
	Payment   domain.Payment     `bson:"payment"`
	Shipment  domain.Shipment    `bson:"shipment"`
	Status    string             `bson:"status"`
	Subtotal  float64            `bson:"subtotal"`
	Total     float64            `bson:"total"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (r *OrderRepositoryMongo) UpdateOrder(ctx context.Context, id string, order domain.Order) error {
	collection := database.GetDB().Collection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": OrderUpdateFields{
			Items:     order.Items,
			Discounts: order.Discounts,
			Taxes:     order.Taxes,
			Payment:   order.Payment,
			Shipment:  order.Shipment,
			Status:    order.Status,
			Subtotal:  order.Subtotal,
			Total:     order.Total,
			UpdatedAt: time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("order not found: %s", id)
	}

	return nil
}

func (r *OrderRepositoryMongo) DeleteOrder(ctx context.Context, id string) error {
	collection := database.GetDB().Collection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
