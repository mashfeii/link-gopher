package storage

import (
	"context"
	"log/slog"
	"sync"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

type InMemoryTagRepository struct {
	tags map[int64]map[string]struct{}
	mu   sync.RWMutex
}

func NewInMemoryTagRepository() *InMemoryTagRepository {
	return &InMemoryTagRepository{
		tags: make(map[int64]map[string]struct{}),
	}
}

func (r *InMemoryTagRepository) AddTagToLink(_ context.Context, linkID int64, tagName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	slog.Info("Adding tag to link", slog.Any("linkID", linkID), slog.Any("tagName", tagName))

	linkTags, ok := r.tags[linkID]
	if !ok {
		r.tags[linkID] = make(map[string]struct{})
	}

	if _, ok := linkTags[tagName]; ok {
		return errors.NewErrTagAlreadyExist(tagName, linkID)
	}

	r.tags[linkID][tagName] = struct{}{}

	return nil
}

func (r *InMemoryTagRepository) RemoveTagsFromLink(_ context.Context, linkID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tags, linkID)

	return nil
}

func (r *InMemoryTagRepository) GetTagsByLink(_ context.Context, linkID int64) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tags := make([]string, 0)

	linkTags, ok := r.tags[linkID]
	if !ok {
		return tags, nil
	}

	for tag := range linkTags {
		tags = append(tags, tag)
	}

	return tags, nil
}
