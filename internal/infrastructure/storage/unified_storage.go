package storage

import (
	"context"
	"slices"
	"sync"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

type CombinedRepository struct {
	users map[int64]*models.User
	mu    sync.RWMutex
}

func NewCombinedRepository() *CombinedRepository {
	return &CombinedRepository{
		users: make(map[int64]*models.User),
		mu:    sync.RWMutex{},
	}
}

func (r *CombinedRepository) CreateUser(_ context.Context, chatID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[chatID]; exists {
		return errors.NewErrUserAlreadyExists(chatID)
	}

	r.users[chatID] = &models.User{
		ChatID:  chatID,
		Links:   make(map[string]models.Link),
		Tags:    make(map[string]struct{}),
		Filters: make(map[string][]string),
	}

	return nil
}

func (r *CombinedRepository) GetUser(_ context.Context, chatID int64) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[chatID]
	if !exists {
		return nil, errors.NewErrUserNotFound()
	}

	return user, nil
}

func (r *CombinedRepository) DeleteUser(_ context.Context, chatID int64) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[chatID]
	if !exists {
		return nil, errors.NewErrUserNotFound()
	}

	delete(r.users, chatID)

	return user, nil
}

func (r *CombinedRepository) AddLink(_ context.Context, chatID int64, link *models.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[chatID]
	if !exists {
		return errors.NewErrUserNotFound()
	}

	if _, exists := user.Links[link.URL]; exists {
		return errors.NewErrLinkAlreadyExists(link.URL)
	}

	user.AddLink(link)

	return nil
}

func (r *CombinedRepository) DeleteLink(_ context.Context, chatID int64, url string) (*models.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[chatID]
	if !exists {
		return nil, errors.NewErrUserNotFound()
	}

	link, exists := user.Links[url]
	if !exists {
		return nil, errors.NewErrLinkNotFound()
	}

	user.DeleteLink(&link)

	return &link, nil
}

func (r *CombinedRepository) GetLinks(_ context.Context, chatID int64) ([]models.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[chatID]
	if !exists {
		return nil, errors.NewErrUserNotFound()
	}

	links := make([]models.Link, 0, len(user.Links))
	for _, link := range user.Links {
		links = append(links, link)
	}

	return links, nil
}

func (r *CombinedRepository) GetAllActiveLinks(_ context.Context) ([]models.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	links := make([]models.Link, 0)

	for _, user := range r.users {
		for _, link := range user.Links {
			links = append(links, link)
		}
	}

	return links, nil
}

func (r *CombinedRepository) AddTagToLink(_ context.Context, chatID int64, url, tag string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[chatID]
	if !exists {
		return errors.NewErrUserNotFound()
	}

	link, exists := user.Links[url]
	if !exists {
		return errors.NewErrLinkNotFound()
	}

	link.AddTag(tag)

	return nil
}

func (r *CombinedRepository) DeleteTagFromLink(_ context.Context, chatID int64, url, tag string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[chatID]
	if !exists {
		return errors.NewErrUserNotFound()
	}

	link, exists := user.Links[url]
	if !exists {
		return errors.NewErrLinkNotFound()
	}

	link.DeleteTag(tag)

	stillUsed := false

	for _, l := range user.Links {
		if _, ok := l.Tags[tag]; ok {
			stillUsed = true
			break
		}
	}

	if !stillUsed {
		delete(user.Tags, tag)
	}

	return nil
}

func (r *CombinedRepository) UpdateLinkFilters(_ context.Context, chatID int64, url string, filters map[string][]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[chatID]
	if !exists {
		return errors.NewErrUserNotFound()
	}

	link, exists := user.Links[url]
	if !exists {
		return errors.NewErrLinkNotFound()
	}

	link.UpdateFilters(filters)

	for key, value := range filters {
		if _, ok := user.Filters[key]; !ok {
			user.Filters[key] = value
			continue
		}

		for _, filter := range value {
			if !slices.Contains(user.Filters[key], filter) {
				user.Filters[key] = append(user.Filters[key], filter)
			}
		}
	}

	return nil
}
