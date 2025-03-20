package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage"
)

func TestInMemoryFilterRepository_AddFilterToLink(t *testing.T) {
	t.Run("Correctly adds filter", func(t *testing.T) {
		repo := storage.NewInMemoryFilterRepository()

		err := repo.AddFilterToLink(context.Background(), 1, "key", "value")
		assert.NoError(t, err)

		filters, _ := repo.GetFiltersByLink(context.Background(), 1)
		assert.Len(t, filters, 1)
		assert.Contains(t, filters, "key:value")
	})

	t.Run("Filter already exists", func(t *testing.T) {
		repo := storage.NewInMemoryFilterRepository()

		err := repo.AddFilterToLink(context.Background(), 1, "key", "value")
		assert.NoError(t, err)

		err = repo.AddFilterToLink(context.Background(), 1, "key", "value")
		assert.ErrorContains(t, err, "already exists")
	})
}

func TestInMemoryFilterRepository_RemoveFiltersFromLink(t *testing.T) {
	t.Run("Correctly removes filter", func(t *testing.T) {
		repo := storage.NewInMemoryFilterRepository()

		err := repo.AddFilterToLink(context.Background(), 1, "key", "value")
		assert.NoError(t, err)

		err = repo.AddFilterToLink(context.Background(), 1, "key2", "value")
		assert.NoError(t, err)

		err = repo.RemoveFiltersFromLink(context.Background(), 1)
		assert.NoError(t, err)

		filters, _ := repo.GetFiltersByLink(context.Background(), 1)
		assert.Empty(t, filters)
		assert.NotNil(t, filters)
	})
}

func TestInMemoryFilterRepository_GetFiltersByLink(t *testing.T) {
	t.Run("Correctly gets filters", func(t *testing.T) {
		repo := storage.NewInMemoryFilterRepository()

		err := repo.AddFilterToLink(context.Background(), 1, "key", "value")
		assert.NoError(t, err)

		err = repo.AddFilterToLink(context.Background(), 1, "key", "value2")
		assert.NoError(t, err)

		err = repo.AddFilterToLink(context.Background(), 1, "key2", "value")
		assert.NoError(t, err)

		filters, err := repo.GetFiltersByLink(context.Background(), 1)
		assert.NoError(t, err)
		assert.Len(t, filters, 3)
	})

	t.Run("Link not found", func(t *testing.T) {
		repo := storage.NewInMemoryFilterRepository()

		filters, err := repo.GetFiltersByLink(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, filters)
		assert.Empty(t, filters)
	})
}
