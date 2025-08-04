package database

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CollectionImpl struct {
	collection *mongo.Collection
}

func NewCollection(collection *mongo.Collection) Collection {
	return &CollectionImpl{collection}
}

func (cl *CollectionImpl) FindOneInto(ctx context.Context, filter interface{}, dest interface{}) error {
	result := cl.collection.FindOne(ctx, filter)
	err := result.Decode(dest)
	if err != nil {
		log.Error().Err(err).Msg("MongoDB FindOneInto failed")
	}
	return err
}

func (cl *CollectionImpl) FindManyInto(ctx context.Context, filter interface{}, dest interface{}, opts ...*options.FindOptions) error {
	cursor, err := cl.collection.Find(ctx, filter, opts...)
	if err != nil {
		log.Error().Err(err).Msg("MongoDB FindManyInto failed: cursor")
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, dest); err != nil {
		log.Error().Err(err).Msg("MongoDB FindManyInto failed: decode")
		return err
	}
	return nil
}

func (cl *CollectionImpl) Insert(ctx context.Context, document interface{}) error {
	_, err := cl.collection.InsertOne(ctx, document)
	if err != nil {
		log.Error().Err(err).Msg("MongoDB Insert failed")
	}
	return err
}

func (cl *CollectionImpl) Update(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	_, err := cl.collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		log.Error().Err(err).Msg("MongoDB Update failed")
	}
	return err
}

func (cl *CollectionImpl) Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) error {
	_, err := cl.collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		log.Error().Err(err).Msg("MongoDB Delete failed")
	}
	return err
}

func (cl *CollectionImpl) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	count, err := cl.collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		log.Error().Err(err).Msg("MongoDB CountDocuments failed")
	}
	return count, err
}
