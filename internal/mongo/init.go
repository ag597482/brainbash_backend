package mongo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"brainbash_backend/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	client   *mongo.Client
	database *mongo.Database
	once     sync.Once
)

// Init initializes the singleton MongoDB client and database.
// Safe to call multiple times; only the first call takes effect.
func Init(cfg *config.AppConfig) {
	once.Do(func() {
		mongoCfg := cfg.StaticConfig.Mongo
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		c, err := mongo.Connect(options.Client().ApplyURI(mongoCfg.URI))
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
		}

		if err := c.Ping(ctx, readpref.Primary()); err != nil {
			panic(fmt.Sprintf("Failed to ping MongoDB: %v", err))
		}

		client = c
		database = c.Database(mongoCfg.Database)

		log.Printf("Connected to MongoDB database: %s", mongoCfg.Database)
	})
}

// GetDatabase returns the singleton MongoDB database instance.
// Panics if Init has not been called.
func GetDatabase() *mongo.Database {
	if database == nil {
		panic("MongoDB not initialized. Call mongo.Init() first.")
	}
	return database
}

// Disconnect gracefully closes the MongoDB connection.
func Disconnect(ctx context.Context) error {
	if client == nil {
		return nil
	}
	return client.Disconnect(ctx)
}
