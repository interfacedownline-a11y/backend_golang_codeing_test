package repository

import (
	"backend_golang_codeing_test/internal/user/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	CountByEmail(ctx context.Context, email string) (int64, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
	GetAllUser(ctx context.Context) ([]*model.User, error)
	UpdateFields(ctx context.Context, id primitive.ObjectID, fields bson.M) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
