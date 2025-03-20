package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage"
)

func TestInMemoryUserRepository_AddUser(t *testing.T) {
	t.Run("Correctly adds a user", func(t *testing.T) {
		repo := storage.NewInMemoryUserRepository()
		user := &models.User{ChatID: 1}

		err := repo.AddUser(context.Background(), user)
		assert.NoError(t, err)

		getUser, err := repo.GetUser(context.Background(), 1)
		assert.NoError(t, err)

		assert.Equal(t, user.ChatID, getUser.ChatID)
	})

	t.Run("User already exists", func(t *testing.T) {
		repo := storage.NewInMemoryUserRepository()
		user := &models.User{ChatID: 1}

		err := repo.AddUser(context.Background(), user)
		assert.NoError(t, err)

		err = repo.AddUser(context.Background(), user)
		assert.ErrorContains(t, err, "already exists")
	})
}

func TestInMemoryUserRepository_GetUser(t *testing.T) {
	t.Run("User not found", func(t *testing.T) {
		repo := storage.NewInMemoryUserRepository()

		getUser, err := repo.GetUser(context.Background(), 1)
		assert.ErrorContains(t, err, "not found")
		assert.Nil(t, getUser)
	})
}

func TestInMemoryUserRepository_DeleteUser(t *testing.T) {
	t.Run("Correct deletion from repository", func(t *testing.T) {
		repo := storage.NewInMemoryUserRepository()
		user := &models.User{ChatID: 1}

		err := repo.AddUser(context.Background(), user)
		assert.NoError(t, err)

		err = repo.DeleteUser(context.Background(), 1)
		assert.NoError(t, err)

		getUser, err := repo.GetUser(context.Background(), 1)
		assert.ErrorContains(t, err, "not found")
		assert.Nil(t, getUser)
	})

	t.Run("User not found", func(t *testing.T) {
		repo := storage.NewInMemoryUserRepository()

		err := repo.DeleteUser(context.Background(), 1)
		assert.ErrorContains(t, err, "not found")
	})
}
