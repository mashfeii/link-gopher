package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/db"
	storage "github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage/sql"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSQLTagRepository(t *testing.T) {
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

	repo := storage.NewSQLTagRepository(pool)
	userRepo := storage.NewSQLUserRepository(pool)
	linksRepo := storage.NewSQLLinkRepository(pool)

	err = userRepo.AddUser(ctx, &models.User{
		ChatID: 123456789,
	})
	if err != nil {
		t.Fatalf("failed to add user: %s", err)
	}

	_, err = linksRepo.AddLink(ctx, &models.Link{
		ChatID:      123456789,
		URL:         "https://example.com",
		Type:        models.LinkTypeGithub,
		LastChecked: time.Now(),
		LastUpdated: time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to add link: %s", err)
	}

	_, err = linksRepo.AddLink(ctx, &models.Link{
		ChatID:      123456789,
		URL:         "https://example.com/2",
		Type:        models.LinkTypeGithub,
		LastChecked: time.Now(),
		LastUpdated: time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to add second link: %s", err)
	}

	t.Run("AddTagToUser", func(t *testing.T) {
		tagName := "test-tag"
		tagID, err := repo.AddTagToUser(ctx, 123456789, tagName)

		assert.NoError(t, err)
		assert.NotNil(t, tagID)
	})

	t.Run("AddTagToLink", func(t *testing.T) {
		tagName := "test-tag"
		err := repo.AddTagToLink(ctx, 123456789, "https://example.com", tagName)

		assert.NoError(t, err)
		err = repo.AddTagToLink(ctx, 123456789, "https://example.com/2", tagName)
		assert.NoError(t, err)
	})

	t.Run("GetTagsByLink", func(t *testing.T) {
		tags, err := repo.GetTagsByLink(ctx, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, tags)
		assert.Contains(t, tags, "test-tag")
	})

	t.Run("GetTagsByUser", func(t *testing.T) {
		tags, err := repo.GetTagsByUser(ctx, 123456789)

		assert.NoError(t, err)
		assert.NotEmpty(t, tags)
		assert.Contains(t, tags, "test-tag")
	})

	t.Run("RemoveTagsFromLink", func(t *testing.T) {
		err := repo.RemoveTagsFromLink(ctx, 1)
		assert.NoError(t, err)

		tags, err := repo.GetTagsByLink(ctx, 1)
		assert.NoError(t, err)
		assert.Empty(t, tags, "Tags should be empty after removal")
	})
}
