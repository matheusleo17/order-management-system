package repository

import (
	"context"
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

	collection := database.DB.Collection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, order)
	return err
}

func (r *OrderRepositoryMongo) GetOrders(ctx context.Context) ([]domain.Order, error) {

	collection := database.DB.Collection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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

func (r *OrderRepositoryMongo) GetOrdersById(ctx context.Context, id string) (domain.Order, error) {

	collection := database.DB.Collection(collectionName)

	var order domain.Order

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (r *OrderRepositoryMongo) UpdateOrder(ctx context.Context, id string, order domain.Order) error {

	collection := database.DB.Collection(collectionName)

	filter := bson.M{"_id": id}

	update := bson.M{
		"$set": order,
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *OrderRepositoryMongo) DeleteOrder(ctx context.Context, id string) error {

	collection := database.DB.Collection(collectionName)

	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
