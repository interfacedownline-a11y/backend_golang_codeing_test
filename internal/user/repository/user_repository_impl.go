package repository

import (
	"backend_golang_codeing_test/internal/user/model"
	"backend_golang_codeing_test/pkg/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryImpl struct {
	userCol database.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	userCol := database.NewCollection(db.Collection("users"))
	return &userRepositoryImpl{userCol: userCol}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	err := r.userCol.Insert(ctx, user)
	return err
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	filter := bson.M{"email": email}
	var user model.User

	err := r.userCol.FindOneInto(ctx, filter, &user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}

	count, err := r.userCol.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.userCol.FindOneInto(ctx, bson.M{"_id": id}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetAllUser(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	err := r.userCol.FindManyInto(ctx, bson.M{}, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepositoryImpl) CountByEmail(ctx context.Context, email string) (int64, error) {
	return r.userCol.CountDocuments(ctx, bson.M{"email": email})
}

func (r *userRepositoryImpl) UpdateFields(ctx context.Context, id primitive.ObjectID, fields bson.M) error {
	update := bson.M{"$set": fields}
	return r.userCol.Update(ctx, bson.M{"_id": id}, update)
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.userCol.Delete(ctx, bson.M{"_id": id})
}

func (r *userRepositoryImpl) CountAll(ctx context.Context) (int64, error) {
	return r.userCol.CountDocuments(ctx, bson.M{})
}
