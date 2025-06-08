package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	scrapper_api "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/servers/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/repository"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

type LinksService interface {
	TrackLink(ctx context.Context, chatID int64, url string, tags, filters []string) (link *models.Link, err error)
	DeleteLink(ctx context.Context, chatID int64, url string) (link *models.Link, meta *LinkMeta, err error)
	RetrieveLinks(ctx context.Context, chatID int64) (links []scrapper_api.LinkResponse, err error)
	RetrieveMeta(ctx context.Context, chatID int64, url string) (meta *LinkMeta, err error)
}

type LinkServiceImpl struct {
	usersRepo   repository.UserRepository
	linksRepo   repository.LinkRepository
	tagsRepo    repository.TagRepository
	filtersRepo repository.FilterRepository
}

func NewLinkService(
	usersRepo repository.UserRepository,
	linksRepo repository.LinkRepository,
	tagsRepo repository.TagRepository,
	filtersRepo repository.FilterRepository,
) *LinkServiceImpl {
	return &LinkServiceImpl{
		usersRepo:   usersRepo,
		linksRepo:   linksRepo,
		tagsRepo:    tagsRepo,
		filtersRepo: filtersRepo,
	}
}

func (s *LinkServiceImpl) TrackLink(
	ctx context.Context,
	chatID int64,
	url string,
	tags, filters []string,
) (*models.Link, error) {
	if _, err := s.usersRepo.GetUser(ctx, chatID); err != nil {
		return nil, err
	}

	link := &models.Link{
		ChatID:      chatID,
		URL:         url,
		LastUpdated: time.Now(),
		LastChecked: time.Now(),
	}

	linkID, err := s.linksRepo.AddLink(ctx, link)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if err := s.tagsRepo.AddTagToLink(ctx, link.LinkID, link.URL, tag); err != nil {
			return nil, err
		}
	}

	for _, filter := range filters {
		keyValue := strings.Split(filter, ":")
		if len(keyValue) != 2 {
			return nil, errors.NewErrInvalidFilterFormat(filter)
		}

		if err := s.filtersRepo.AddFilterToLink(ctx, *linkID, keyValue[0], keyValue[1]); err != nil {
			return nil, err
		}
	}

	return link, nil
}

func (s *LinkServiceImpl) DeleteLink(ctx context.Context, chatID int64, url string) (*models.Link, *LinkMeta, error) {
	if _, err := s.usersRepo.GetUser(ctx, chatID); err != nil {
		return nil, nil, err
	}

	meta, err := s.RetrieveMeta(ctx, chatID, url)
	if err != nil {
		return nil, nil, err
	}

	link, err := s.linksRepo.DeleteLink(ctx, chatID, url)
	if err != nil {
		return nil, nil, err
	}

	if err := s.filtersRepo.RemoveFiltersFromLink(ctx, link.LinkID); err != nil {
		return nil, nil, err
	}

	if err := s.tagsRepo.RemoveTagsFromLink(ctx, link.LinkID); err != nil {
		return nil, nil, err
	}

	return link, meta, nil
}

func (s *LinkServiceImpl) RetrieveLinks(ctx context.Context, chatID int64) ([]scrapper_api.LinkResponse, error) {
	const op = "links_service.RetrieveLinks"

	if _, err := s.usersRepo.GetUser(ctx, chatID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, errors.NewErrUserNotFound(chatID))
	}

	userLinks, err := s.linksRepo.GetLinksByUser(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	linksResponse := make([]scrapper_api.LinkResponse, 0, len(userLinks))

	for _, link := range userLinks {
		meta, err := s.RetrieveMeta(ctx, chatID, link.URL)
		if err != nil {
			return nil, fmt.Errorf("%s: unable to retrieve meta for link %s: %w", op, link.URL, err)
		}

		linksResponse = append(linksResponse, scrapper_api.LinkResponse{
			Id:      &link.LinkID,
			Url:     &link.URL,
			Filters: &meta.Filters,
			Tags:    &meta.Tags,
		})
	}

	return linksResponse, nil
}

func (s *LinkServiceImpl) RetrieveMeta(ctx context.Context, chatID int64, url string) (*LinkMeta, error) {
	links, err := s.linksRepo.GetLinksByUser(ctx, chatID)
	if err != nil {
		return nil, err
	}

	var link *models.Link

	for _, l := range links {
		if strings.EqualFold(l.URL, url) {
			link = &l
			break
		}
	}

	filters, err := s.filtersRepo.GetFiltersByLink(ctx, link.LinkID)
	if err != nil {
		return nil, err
	}

	tags, err := s.tagsRepo.GetTagsByLink(ctx, link.LinkID)
	if err != nil {
		return nil, err
	}

	return &LinkMeta{Tags: tags, Filters: filters}, nil
}
