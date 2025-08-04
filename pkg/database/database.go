package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	GetDB() *mongo.Database
}

type Collection interface {
	FindOneInto(ctx context.Context, filter interface{}, dest interface{}) error
	FindManyInto(ctx context.Context, filter interface{}, dest interface{}, opts ...*options.FindOptions) error
	Insert(ctx context.Context, document interface{}) error
	Update(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error
	Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) error
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
}
