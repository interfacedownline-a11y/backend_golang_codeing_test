package service

import (
	"backend_golang_codeing_test/internal/user/model"
	"backend_golang_codeing_test/internal/user/repository"
	"backend_golang_codeing_test/internal/user/transport/http/dto"
	"context"
	"fmt"
	"net/mail"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) GetAllUser(ctx context.Context) ([]*model.User, error) {
	return s.repo.GetAllUser(ctx)
}

func (s *userServiceImpl) FindByID(ctx context.Context, id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object id: %w", err)
	}

	user, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return user, nil
}

func (s *userServiceImpl) Update(ctx context.Context, req dto.UpdateUserRequest) error {
	objectID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return fmt.Errorf("invalid object id: %w", err)
	}

	user, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	updateFields := bson.M{}

	if req.Email != "" && req.Email != user.Email {
		if _, err := mail.ParseAddress(req.Email); err != nil {
			return fmt.Errorf("invalid email format: %w", err)
		}

		count, err := s.repo.CountByEmail(ctx, req.Email)
		if err != nil {
			return fmt.Errorf("failed to check email existence: %w", err)
		}
		if count > 0 {
			return fmt.Errorf("email already exists")
		}

		updateFields["email"] = req.Email
	}

	if req.Name != "" && req.Name != user.Name {
		updateFields["name"] = req.Name
	}

	if len(updateFields) == 0 {
		return nil
	}

	return s.repo.UpdateFields(ctx, objectID, updateFields)
}

func (s *userServiceImpl) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object id: %w", err)
	}

	if err := s.repo.Delete(ctx, objectID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
