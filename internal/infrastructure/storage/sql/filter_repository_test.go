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

func TestFilterRepository(t *testing.T) {
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

	repo := storage.NewSQLFilterRepository(pool)
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

	t.Run("AddFilterToLink", func(t *testing.T) {
		filterKey := "test"
		filterValue := "value"

		err := repo.AddFilterToLink(ctx, 1, filterKey, filterValue)
		assert.NoError(t, err)
		err = repo.AddFilterToLink(ctx, 2, filterKey, filterValue)
		assert.NoError(t, err)
	})

	t.Run("GetFiltersByLink", func(t *testing.T) {
		filters, err := repo.GetFiltersByLink(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, filters, 1)
		assert.Equal(t, "test:value", filters[0])
	})

	t.Run("RemoveFiltersFromLink", func(t *testing.T) {
		err := repo.RemoveFiltersFromLink(ctx, 1)
		assert.NoError(t, err)
		filters, err := repo.GetFiltersByLink(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, filters, 0)
	})
}
