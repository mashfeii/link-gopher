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

func TestSQLUpdateRepository(t *testing.T) {
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

	repo := storage.NewSQLUpdateRepository(pool)
	repoUsers := storage.NewSQLUserRepository(pool)
	repoLinks := storage.NewSQLLinkRepository(pool)

	user := &models.User{
		ChatID: 1,
	}

	err = repoUsers.AddUser(ctx, user)
	if err != nil {
		t.Fatalf("failed to add user: %s", err)
	}

	linkForUpdate := &models.Link{
		ChatID:      1,
		URL:         "https://example.com",
		Type:        models.LinkTypeGithub,
		LastChecked: time.Now(),
		LastUpdated: time.Now(),
	}

	_, err = repoLinks.AddLink(ctx, linkForUpdate)
	if err != nil {
		t.Fatalf("failed to add link for update: %s", err)
	}

	t.Run("SaveUpdate", func(t *testing.T) {
		update := &models.Update{
			LinkID:      1,
			UpdateType:  models.UpdateTypePullRequest,
			Title:       "Test Update",
			UserName:    "testuser",
			BodyPreview: "This is a test update.",
			CreatedAt:   time.Now(),
			Sent:        false,
		}

		err = repo.SaveUpdate(ctx, update)
		assert.NoError(t, err, "SaveUpdate should not return an error")

		update2 := &models.Update{
			LinkID:      1,
			UpdateType:  models.UpdateTypeIssue,
			Title:       "Test Update 2",
			UserName:    "testuser2",
			BodyPreview: "This is another test update.",
			CreatedAt:   time.Now(),
			Sent:        false,
		}

		err = repo.SaveUpdate(ctx, update2)
		assert.NoError(t, err, "SaveUpdate should not return an error for second update")
	})

	t.Run("GetUpdateByLink", func(t *testing.T) {
		updates, err := repo.GetUpdateByLink(ctx, 1, 10, 0)

		assert.NoError(t, err, "GetUpdateByLink should not return an error")
		assert.Len(t, updates, 2, "GetUpdateByLink should return 2 updates")
	})
}
