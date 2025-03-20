package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/storage"
)

func TestInMemoryTagRepository_AddTagToLink(t *testing.T) {
	t.Run("Correctly adds a tag to a link", func(t *testing.T) {
		tagRepo := storage.NewInMemoryTagRepository()

		err := tagRepo.AddTagToLink(context.Background(), 1, "tag1")
		assert.NoError(t, err)

		err = tagRepo.AddTagToLink(context.Background(), 1, "tag2")
		assert.NoError(t, err)

		tags, err := tagRepo.GetTagsByLink(context.Background(), 1)
		assert.NoError(t, err)
		assert.Contains(t, tags, "tag1")
		assert.Contains(t, tags, "tag2")
	})

	t.Run("Tag already exists", func(t *testing.T) {
		tagRepo := storage.NewInMemoryTagRepository()

		err := tagRepo.AddTagToLink(context.Background(), 1, "tag1")
		assert.NoError(t, err)

		err = tagRepo.AddTagToLink(context.Background(), 1, "tag1")
		assert.ErrorContains(t, err, "already exists")
	})
}

func TestInMemoryTagRepository_RemoveTagsFromLink(t *testing.T) {
	t.Run("Correctly removes tags from a link", func(t *testing.T) {
		tagRepo := storage.NewInMemoryTagRepository()

		err := tagRepo.AddTagToLink(context.Background(), 1, "tag1")
		assert.NoError(t, err)

		err = tagRepo.AddTagToLink(context.Background(), 1, "tag2")
		assert.NoError(t, err)

		err = tagRepo.RemoveTagsFromLink(context.Background(), 1)
		assert.NoError(t, err)

		tags, err := tagRepo.GetTagsByLink(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, tags)
		assert.Empty(t, tags)
	})
}

func TestInMemoryTagRepository_GetTagsByLink(t *testing.T) {
	t.Run("Correctly returns tags for a link", func(t *testing.T) {
		tagRepo := storage.NewInMemoryTagRepository()

		err := tagRepo.AddTagToLink(context.Background(), 1, "tag1")
		assert.NoError(t, err)

		err = tagRepo.AddTagToLink(context.Background(), 1, "tag2")
		assert.NoError(t, err)

		tags, err := tagRepo.GetTagsByLink(context.Background(), 1)
		assert.NoError(t, err)
		assert.Contains(t, tags, "tag1")
		assert.Contains(t, tags, "tag2")
	})

	t.Run("Link not found", func(t *testing.T) {
		tagRepo := storage.NewInMemoryTagRepository()

		tags, err := tagRepo.GetTagsByLink(context.Background(), 1)
		assert.NoError(t, err)
		assert.Empty(t, tags)
	})
}
