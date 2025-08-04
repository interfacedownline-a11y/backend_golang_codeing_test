package service

import (
	"backend_golang_codeing_test/internal/auth/model"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, input model.RegisterRequest) error
	Login(ctx context.Context, input model.LoginRequest) (*model.LoginResponse, error)
}
