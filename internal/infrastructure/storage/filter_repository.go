package storage

import (
	"context"
	"slices"
	"sync"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

type InMemoryFilterRepository struct {
	filters map[int64]map[string][]string
	mu      sync.RWMutex
}

func NewInMemoryFilterRepository() *InMemoryFilterRepository {
	return &InMemoryFilterRepository{
		filters: make(map[int64]map[string][]string),
	}
}

func (r *InMemoryFilterRepository) AddFilterToLink(_ context.Context, linkID int64, key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.filters[linkID]; !ok {
		r.filters[linkID] = make(map[string][]string)
	}

	if slices.Contains(r.filters[linkID][key], value) {
		return errors.NewErrFilterAlreadyExists(linkID, key, value)
	}

	r.filters[linkID][key] = append(r.filters[linkID][key], value)

	return nil
}

func (r *InMemoryFilterRepository) RemoveFiltersFromLink(_ context.Context, linkID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.filters, linkID)

	return nil
}

func (r *InMemoryFilterRepository) GetFiltersByLink(_ context.Context, linkID int64) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filters := make([]string, 0)

	if _, ok := r.filters[linkID]; !ok {
		return filters, nil
	}

	for key, values := range r.filters[linkID] {
		for _, value := range values {
			filters = append(filters, key+":"+value)
		}
	}

	return filters, nil
}
