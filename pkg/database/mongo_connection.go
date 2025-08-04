package database

import (
	"backend_golang_codeing_test/config"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	db     *mongo.Database
	once   sync.Once
)

func NewMongoDatabase(conf *config.Database) *mongo.Database {
	once.Do(func() {
		uri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
			conf.User,
			conf.Password,
			conf.Host,
			conf.Port,
		)

		fmt.Println("MongoDB URI:", uri)
		clientOptions := options.Client().ApplyURI(uri)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("failed to connect to MongoDB: %v", err)
		}

		if err := client.Ping(ctx, nil); err != nil {
			log.Fatalf("MongoDB ping failed: %v", err)
		}

		db = client.Database(conf.DBName)
		fmt.Println("Connected to MongoDB")
	})
	return db
}
