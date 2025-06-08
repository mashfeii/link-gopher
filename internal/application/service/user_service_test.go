package service_test

import (
	"context"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/application/service"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterUser(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	userService := service.NewUserService(mockRepo)

	chatID := int64(12345)
	mockRepo.On("AddUser", mock.Anything, &models.User{ChatID: chatID}).Return(nil).Once()

	err := userService.RegisterUser(context.Background(), chatID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	userService := service.NewUserService(mockRepo)

	chatID := int64(12345)
	expectedUser := &models.User{ChatID: chatID}
	mockRepo.On("GetUser", mock.Anything, chatID).Return(expectedUser, nil).Once()

	user, err := userService.GetUser(context.Background(), chatID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	userService := service.NewUserService(mockRepo)

	chatID := int64(12345)
	mockRepo.On("GetUser", mock.Anything, chatID).Return(nil, errors.ErrUserNotFound{}).Once()

	user, err := userService.GetUser(context.Background(), chatID)
	assert.Nil(t, user)
	assert.Equal(t, errors.ErrUserNotFound{}, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_Error(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	userService := service.NewUserService(mockRepo)

	chatID := int64(12345)
	mockRepo.On("GetUser", mock.Anything, chatID).Return(nil, assert.AnError).Once()

	user, err := userService.GetUser(context.Background(), chatID)
	assert.Nil(t, user)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	userService := service.NewUserService(mockRepo)

	chatID := int64(12345)
	mockRepo.On("DeleteUser", mock.Anything, chatID).Return(nil).Once()

	err := userService.DeleteUser(context.Background(), chatID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock(t)
	userService := service.NewUserService(mockRepo)

	chatID := int64(12345)
	mockRepo.On("DeleteUser", mock.Anything, chatID).Return(errors.ErrUserNotFound{}).Once()
	err := userService.DeleteUser(context.Background(), chatID)
	assert.Equal(t, errors.ErrUserNotFound{}, err)
	mockRepo.AssertExpectations(t)
}
