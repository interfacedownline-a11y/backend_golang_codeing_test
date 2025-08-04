package service

import (
	"backend_golang_codeing_test/internal/user/mocks"
	"backend_golang_codeing_test/internal/user/model"
	"backend_golang_codeing_test/internal/user/transport/http/dto"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindByID_ValidID_ReturnsUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := NewUserService(mockRepo)

	id := primitive.NewObjectID()
	expected := &model.User{ID: id, Name: "Alice"}

	mockRepo.On("FindByID", mock.Anything, id).Return(expected, nil)

	user, err := svc.FindByID(context.TODO(), id.Hex())

	require.NoError(t, err)
	assert.Equal(t, expected.Name, user.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUser_ReturnsUsers(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := NewUserService(mockRepo)

	expected := []*model.User{
		{ID: primitive.NewObjectID(), Name: "Alice"},
		{ID: primitive.NewObjectID(), Name: "Bob"},
	}

	mockRepo.On("GetAllUser", mock.Anything).Return(expected, nil)

	users, err := svc.GetAllUser(context.TODO())

	require.NoError(t, err)
	assert.Equal(t, len(expected), len(users))
	mockRepo.AssertExpectations(t)
}

func TestUpdate_ValidInput_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := NewUserService(mockRepo)

	id := primitive.NewObjectID()
	req := dto.UpdateUserRequest{
		ID:    id.Hex(),
		Name:  "Alice Updated",
		Email: "alice@example.com",
	}
	user := &model.User{ID: id, Name: "Alice", Email: "old@example.com"}

	mockRepo.On("FindByID", mock.Anything, id).Return(user, nil)
	mockRepo.On("CountByEmail", mock.Anything, req.Email).Return(int64(0), nil)
	mockRepo.On("UpdateFields", mock.Anything, id, mock.Anything).Return(nil)

	err := svc.Update(context.TODO(), req)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete_ValidID_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := NewUserService(mockRepo)

	id := primitive.NewObjectID()

	mockRepo.On("Delete", mock.Anything, id).Return(nil)

	err := svc.Delete(context.TODO(), id.Hex())

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
