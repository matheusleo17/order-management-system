package repository

import (
	"context"
	"time"

	"order-management-system/internal/database"
	"order-management-system/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
)

const collectionName = "orders"

func InsertOrder(order domain.Order) error {

	collection := database.DB.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, order)

	return err
}

func GetOrders() ([]domain.Order, error) {

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
