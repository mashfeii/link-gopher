package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage"
)

func TestInMemoryLinkRepository_AddLink(t *testing.T) {
	t.Run("Link already exists", func(t *testing.T) {
		repo := storage.NewInMemoryLinkRepository()
		link := models.Link{
			ChatID:     1,
			URL:        "https://example.com",
			LastUpdate: time.Now(),
		}

		_, err := repo.AddLink(context.Background(), &link)
		assert.NoError(t, err)

		_, err = repo.AddLink(context.Background(), &link)
		assert.ErrorContains(t, err, "already exists")
	})

	t.Run("Correctly adds a link", func(t *testing.T) {
		repo := storage.NewInMemoryLinkRepository()
		link := models.Link{
			ChatID:     1,
			URL:        "https://example.com",
			LastUpdate: time.Now(),
		}

		id, err := repo.AddLink(context.Background(), &link)
		assert.NoError(t, err)

		retrievedLink, err := repo.GetLinkByURL(context.Background(), 1, "https://example.com")
		assert.NoError(t, err)

		assert.Equal(t, *id, retrievedLink.LinkID)
	})
}

func TestInMemoryLinkRepository_DeleteLink(t *testing.T) {
	t.Run("Link not found", func(t *testing.T) {
		repo := storage.NewInMemoryLinkRepository()
		link := models.Link{
			ChatID:     1,
			URL:        "https://example.com",
			LastUpdate: time.Now(),
		}

		_, err := repo.AddLink(context.Background(), &link)
		assert.NoError(t, err)

		_, err = repo.DeleteLink(context.Background(), 2, "https://example.com")
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("Correctly deletes a link", func(t *testing.T) {
		repo := storage.NewInMemoryLinkRepository()
		link := models.Link{
			ChatID:     1,
			URL:        "https://example.com",
			LastUpdate: time.Now(),
		}

		_, err := repo.AddLink(context.Background(), &link)
		assert.NoError(t, err)

		deletedLink, err := repo.DeleteLink(context.Background(), 1, "https://example.com")
		assert.NoError(t, err)
		assert.Equal(t, link, *deletedLink)
	})
}

func TestInMemoryLinkRepository_GetLinksByUser(t *testing.T) {
	t.Run("Empty links", func(t *testing.T) {
		repo := storage.NewInMemoryLinkRepository()
		links, err := repo.GetLinksByUser(context.Background(), 1)
		assert.NoError(t, err)
		assert.Empty(t, links)
	})

	t.Run("Get links by user", func(t *testing.T) {
		repo := storage.NewInMemoryLinkRepository()

		link1 := models.Link{
			ChatID: 1,
			URL:    "https://example.com",
		}
		link2 := models.Link{
			ChatID: 1,
			URL:    "https://example2.com",
		}

		_, err := repo.AddLink(context.Background(), &link1)
		assert.NoError(t, err)

		_, err = repo.AddLink(context.Background(), &link2)
		assert.NoError(t, err)

		links, err := repo.GetLinksByUser(context.Background(), 1)
		assert.NoError(t, err)
		assert.Len(t, links, 2)
	})
}

func TestInMemoryLinkRepository_GetLinkByURL(t *testing.T) {
	t.Run("Link not found", func(t *testing.T) {
		repo := storage.NewInMemoryLinkRepository()
		link, err := repo.GetLinkByURL(context.Background(), 1, "https://example.com")
		assert.ErrorContains(t, err, "not found")
		assert.Nil(t, link)
	})

	t.Run("Get link by URL", func(t *testing.T) {
		repo := storage.NewInMemoryLinkRepository()
		link := models.Link{
			ChatID: 1,
			URL:    "https://example.com",
		}

		_, err := repo.AddLink(context.Background(), &link)
		assert.NoError(t, err)

		foundLink, err := repo.GetLinkByURL(context.Background(), 1, "https://example.com")
		assert.NoError(t, err)
		assert.Equal(t, link, *foundLink)
	})
}

func TestInMemoryLinkRepository_GetAllActiveLinks(t *testing.T) {
	t.Run("Empty links", func(t *testing.T) {
		repo := storage.NewInMemoryLinkRepository()
		links, err := repo.GetAllActiveLinks(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, links)
	})

	t.Run("Get all active links", func(t *testing.T) {
		repo := storage.NewInMemoryLinkRepository()

		link1 := models.Link{
			ChatID: 1,
			URL:    "https://example.com",
		}
		link2 := models.Link{
			ChatID: 1,
			URL:    "https://example2.com",
		}

		_, err := repo.AddLink(context.Background(), &link1)
		assert.NoError(t, err)

		_, err = repo.AddLink(context.Background(), &link2)
		assert.NoError(t, err)

		links, err := repo.GetAllActiveLinks(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 2, len(links))
	})
}
