package service

import (
	"backend_golang_codeing_test/internal/user/model"
	"backend_golang_codeing_test/internal/user/transport/http/dto"
	"context"
)

type UserService interface {
	GetAllUser(ctx context.Context) ([]*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user dto.UpdateUserRequest) error
	Delete(ctx context.Context, id string) error
}
