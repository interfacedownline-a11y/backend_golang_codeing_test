package migrations

import (
	"context"
	"log"
	"time"

	"backend_golang_codeing_test/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollections struct {
	User database.Collection
}

func AutoMigrate(db *mongo.Database) {
	SetupCollection(db.Collection("users"))
}

var Collections *MongoCollections

func SetupCollection(db *mongo.Collection) {
	Collections = &MongoCollections{
		User: database.NewCollection(NewUserCollection(db.Database().Collection("users"))),
	}
}

func NewUserCollection(cl *mongo.Collection) *mongo.Collection {
	if err := setupUserCollection(cl); err != nil {
		log.Fatalf("Failed to set up collection: %v", err)
	}
	return cl
}

func setupUserCollection(cl *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "createdAt", Value: 1}},
			Options: options.Index().SetName("created_at_index"),
		},
		{
			Keys:    bson.D{{Key: "name", Value: 1}},
			Options: options.Index().SetName("name_index"),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetName("email_index").SetUnique(true),
		},
	}

	_, err := cl.Indexes().CreateMany(ctx, indexes)
	return err
}
