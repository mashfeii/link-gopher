package storage_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/db"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage/sql"
)

func TestSQLUserRepository(t *testing.T) {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx, "postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	defer func() {
		if err = testcontainers.TerminateContainer(container); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container's host: %s", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container's port: %s", err)
	}

	testConfig := &config.Database{
		Host:     host,
		Port:     port.Int(),
		Username: "postgres",
		Password: "postgres",
		Name:     "testdb",
	}

	pool, err := db.NewPool(testConfig)
	if err != nil {
		t.Fatalf("failed to create database pool: %s", err)
	}
	defer pool.Close()

	repo := storage.NewSQLUserRepository(pool)

	t.Run("AddUser", func(t *testing.T) {
		user := &models.User{ChatID: 12345}
		err := repo.AddUser(ctx, user)
		assert.NoError(t, err)

		var chatID int64
		err = pool.QueryRow(ctx, "SELECT chat_id FROM users WHERE chat_id = $1", user.ChatID).Scan(&chatID)
		assert.NoError(t, err)
		assert.Equal(t, user.ChatID, chatID)
	})

	t.Run("GetUser", func(t *testing.T) {
		user, err := repo.GetUser(ctx, 12345)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(12345), user.ChatID)
	})

	t.Run("GetUserNotFound", func(t *testing.T) {
		user, err := repo.GetUser(ctx, 99999)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		err := repo.DeleteUser(ctx, 12345)
		assert.NoError(t, err)

		var chatID int64
		err = pool.QueryRow(ctx, "SELECT chat_id FROM users WHERE chat_id = $1", 12345).Scan(&chatID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, pgx.ErrNoRows))
	})
}
