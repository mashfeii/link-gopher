package storage

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

type InMemoryLinkRepository struct {
	links     map[int64]*models.Link
	userLinks map[int64][]int64
	urlIndex  map[string]int64
	mu        sync.RWMutex
	seq       int64
}

func NewInMemoryLinkRepository() *InMemoryLinkRepository {
	return &InMemoryLinkRepository{
		links:     make(map[int64]*models.Link),
		userLinks: make(map[int64][]int64),
		urlIndex:  make(map[string]int64),
		seq:       1,
	}
}

func (r *InMemoryLinkRepository) AddLink(_ context.Context, link *models.Link) (*int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%d:%s", link.ChatID, link.URL)
	if _, exists := r.urlIndex[key]; exists {
		return nil, errors.NewErrLinkAlreadyExists(link.URL)
	}

	link.LinkID = r.seq
	r.links[r.seq] = link
	r.userLinks[link.ChatID] = append(r.userLinks[link.ChatID], r.seq)
	r.urlIndex[key] = r.seq
	r.seq++

	return &link.LinkID, nil
}

func (r *InMemoryLinkRepository) DeleteLink(_ context.Context, chatID int64, url string) (*models.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%d:%s", chatID, url)

	linkID, exists := r.urlIndex[key]
	if !exists {
		return nil, errors.NewErrLinkNotFound()
	}

	link := r.links[linkID]
	delete(r.links, linkID)
	delete(r.urlIndex, key)

	links := r.userLinks[chatID]
	for i, id := range links {
		if linkID == id {
			r.userLinks[chatID] = slices.Delete(links, i, i+1)
			break
		}
	}

	return link, nil
}

func (r *InMemoryLinkRepository) GetLinksByUser(_ context.Context, chatID int64) ([]models.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	links := make([]models.Link, 0)

	for _, linkID := range r.userLinks[chatID] {
		links = append(links, *r.links[linkID])
	}

	return links, nil
}

func (r *InMemoryLinkRepository) GetLinkByURL(_ context.Context, chatID int64, url string) (*models.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := fmt.Sprintf("%d:%s", chatID, url)

	linkID, exists := r.urlIndex[key]
	if !exists {
		return nil, errors.NewErrLinkNotFound()
	}

	return r.links[linkID], nil
}

func (r *InMemoryLinkRepository) GetAllActiveLinks(_ context.Context) ([]models.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	links := make([]models.Link, 0, len(r.links))
	for _, link := range r.links {
		links = append(links, *link)
	}

	return links, nil
}
