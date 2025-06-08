package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/db"
	storage "github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage/orm"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSQLLinksRepository(t *testing.T) {
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

	repo := storage.NewORMLinksRepository(pool)
	userRepo := storage.NewORMUserRepository(pool)

	users := []*models.User{
		{ChatID: 12345},
		{ChatID: 67890},
	}

	for _, user := range users {
		err := userRepo.AddUser(ctx, user)
		if err != nil {
			t.Fatalf("failed to add user: %v", err)
		}
	}

	t.Run("AddLink", func(t *testing.T) {
		link := &models.Link{
			ChatID:      12345,
			URL:         "https://example.com",
			Type:        "example",
			LastChecked: time.Now(),
			LastUpdated: time.Now(),
		}

		linkID, err := repo.AddLink(ctx, link)
		if err != nil {
			t.Fatalf("failed to add link: %v", err)
		}

		assert.NotNil(t, linkID, "expected link ID to be non-nil")
		assert.Greater(t, *linkID, int64(0), "expected link ID to be greater than 0")
	})

	t.Run("DeleteLink", func(t *testing.T) {
		deletedLink, err := repo.DeleteLink(ctx, 12345, "https://example.com")
		if err != nil {
			t.Fatalf("failed to delete link: %v", err)
		}

		assert.NotNil(t, deletedLink, "expected deleted link to be non-nil")
		assert.Equal(t, int64(12345), deletedLink.ChatID, "expected deleted link chat ID to match")
		assert.Equal(t, "https://example.com", deletedLink.URL, "expected deleted link URL to match")
		assert.Equal(t, models.LinkType("example"), deletedLink.Type, "expected deleted link type to match")
	})

	t.Run("GetEmptyLinksByUser", func(t *testing.T) {
		links, err := repo.GetLinksByUser(ctx, 12345)

		assert.NoError(t, err, "expected no error when getting links by user")
		assert.Empty(t, links, "expected no links for user with no links")
	})

	t.Run("AddMultipleLinks", func(t *testing.T) {
		links := []*models.Link{
			{ChatID: 12345, URL: "https://example1.com", Type: "example", LastChecked: time.Now(), LastUpdated: time.Now()},
			{ChatID: 12345, URL: "https://example2.com", Type: "example", LastChecked: time.Now(), LastUpdated: time.Now()},
			{ChatID: 12345, URL: "https://example3.com", Type: "example", LastChecked: time.Now(), LastUpdated: time.Now()},
			{ChatID: 12345, URL: "https://example4.com", Type: "example", LastChecked: time.Now(), LastUpdated: time.Now()},
		}

		for _, link := range links {
			linkID, err := repo.AddLink(ctx, link)
			assert.NoError(t, err, "expected no error when adding link")
			assert.NotNil(t, linkID, "expected link ID to be non-nil")
			assert.Greater(t, *linkID, int64(0), "expected link ID to be greater than 0")
		}
	})

	t.Run("GetLinksByUser", func(t *testing.T) {
		links, err := repo.GetLinksByUser(ctx, 12345)
		assert.NoError(t, err, "expected no error when getting links by user")
		assert.Len(t, links, 4, "expected 4 links for user with multiple links")
		assert.Equal(t, "https://example1.com", links[0].URL, "expected first link URL to match")
		assert.Equal(t, "https://example2.com", links[1].URL, "expected second link URL to match")
		assert.Equal(t, "https://example3.com", links[2].URL, "expected third link URL to match")
		assert.Equal(t, "https://example4.com", links[3].URL, "expected fourth link URL to match")
	})

	t.Run("GetLinkByURL", func(t *testing.T) {
		link, err := repo.GetLinkByURL(ctx, 12345, "https://example1.com")
		assert.NoError(t, err, "expected no error when getting link by URL")
		assert.NotNil(t, link, "expected link to be non-nil")
		assert.Equal(t, models.LinkType("example"), link.Type, "expected link type to match")
	})

	t.Run("GetAllActiveLinks", func(t *testing.T) {
		link := &models.Link{
			ChatID:      67890,
			URL:         "https://active-example.com",
			Type:        "example",
			LastChecked: time.Now(),
			LastUpdated: time.Now(),
		}

		_, err := repo.AddLink(ctx, link)
		assert.NoError(t, err, "expected no error when adding active link")

		activeLinks, err := repo.GetAllActiveLinks(ctx)

		assert.NoError(t, err, "expected no error when getting all active links")
		assert.NotEmpty(t, activeLinks, "expected active links to be non-empty")
		assert.Equal(t, 5, len(activeLinks), "expected 5 active links including the newly added one")
	})

	t.Run("UpdateLink", func(t *testing.T) {
		existingLink, err := repo.GetLinkByURL(ctx, 12345, "https://example1.com")
		assert.NoError(t, err)
		assert.NotNil(t, existingLink)

		newLastChecked := time.Now().UTC().Truncate(time.Microsecond)
		newLastUpdated := newLastChecked.Add(24 * time.Hour).Truncate(time.Microsecond)

		link := &models.Link{
			ChatID:      12345,
			URL:         "https://example1.com",
			LastChecked: newLastChecked,
			LastUpdated: newLastUpdated,
		}

		err = repo.UpdateLink(ctx, link)
		assert.NoError(t, err, "expected no error when updating link")

		updatedLink, err := repo.GetLinkByURL(ctx, 12345, "https://example1.com")
		assert.NoError(t, err, "expected no error when getting updated link by URL")
		assert.NotNil(t, updatedLink, "expected updated link to be non-nil")

		assert.WithinDuration(t, link.LastChecked, updatedLink.LastChecked, time.Microsecond, "expected last checked time to match")
		assert.WithinDuration(t, link.LastUpdated, updatedLink.LastUpdated, time.Microsecond, "expected last updated time to match")

		assert.Equal(t, models.LinkType("example"), updatedLink.Type, "expected type to remain unchanged")
	})
}
