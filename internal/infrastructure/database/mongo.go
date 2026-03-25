package database

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	db   *mongo.Database
	once sync.Once
)

func ConnectMongo(log *zap.Logger) error {
	var connectErr error

	once.Do(func() {
		uri := os.Getenv("MONGO_URI")
		if uri == "" {
			uri = "mongodb://localhost:27017"
		}

		dbName := os.Getenv("MONGO_DB")
		if dbName == "" {
			dbName = "order_management"
		}

		client, err := mongo.NewClient(options.Client().ApplyURI(uri))
		if err != nil {
			connectErr = fmt.Errorf("failed to create mongo client: %w", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err = client.Connect(ctx); err != nil {
			connectErr = fmt.Errorf("failed to connect to mongo: %w", err)
			return
		}

		if err = client.Ping(ctx, nil); err != nil {
			connectErr = fmt.Errorf("failed to ping mongo: %w", err)
			return
		}

		db = client.Database(dbName)
	})

	return connectErr
}

func GetDB() *mongo.Database {
	return db
}
