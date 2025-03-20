package service

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	scrapper_api "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/servers/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/repository"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

type UserService interface {
	RegisterUser(ctx context.Context, chatID int64) (status int, err error)
	DeleteUser(ctx context.Context, chatID int64) (status int, err error)
	TrackLink(ctx context.Context, chatID int64, url string, tags, filters []string) (link *models.Link, status int, err error)
	DeleteLink(ctx context.Context, chatID int64, url string) (link *models.Link, meta *LinkMeta, status int, err error)
	RetrieveLinks(ctx context.Context, chatID int64) (links []scrapper_api.LinkResponse, status int, err error)
	RetrieveMeta(ctx context.Context, chatID int64, url string) (meta *LinkMeta, code int, err error)
}

type LinkMeta struct {
	Tags    []string
	Filters []string
}

type Service struct {
	usersRepo   repository.UserRepository
	linksRepo   repository.LinkRepository
	tagsRepo    repository.TagRepository
	filtersRepo repository.FilterRepository
}

func NewService(
	usersRepo repository.UserRepository,
	linksRepo repository.LinkRepository,
	tagsRepo repository.TagRepository,
	filtersRepo repository.FilterRepository,
) *Service {
	return &Service{
		usersRepo:   usersRepo,
		linksRepo:   linksRepo,
		tagsRepo:    tagsRepo,
		filtersRepo: filtersRepo,
	}
}

func (s *Service) RegisterUser(ctx context.Context, chatID int64) (status int, err error) {
	if err := s.usersRepo.AddUser(ctx, &models.User{ChatID: chatID}); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func (s *Service) DeleteUser(ctx context.Context, chatID int64) (status int, err error) {
	if err := s.usersRepo.DeleteUser(ctx, chatID); err != nil {
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}

func (s *Service) TrackLink(
	ctx context.Context,
	chatID int64,
	url string,
	tags, filters []string,
) (
	link *models.Link,
	status int,
	err error,
) {
	if _, err := s.usersRepo.GetUser(ctx, chatID); err != nil {
		return nil, http.StatusUnauthorized, err
	}

	link = &models.Link{
		ChatID:     chatID,
		URL:        url,
		LastUpdate: time.Now(),
	}

	linkID, err := s.linksRepo.AddLink(ctx, link)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	for _, tag := range tags {
		if err := s.tagsRepo.AddTagToLink(ctx, *linkID, tag); err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	for _, filter := range filters {
		keyValue := strings.Split(filter, ":")
		if len(keyValue) != 2 {
			return nil, http.StatusBadRequest, errors.NewErrInvalidFilterFormat(filter)
		}

		if err := s.filtersRepo.AddFilterToLink(ctx, *linkID, keyValue[0], keyValue[1]); err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	return link, http.StatusOK, nil
}

func (s *Service) DeleteLink(ctx context.Context, chatID int64, url string) (link *models.Link, meta *LinkMeta, status int, err error) {
	if _, err := s.usersRepo.GetUser(ctx, chatID); err != nil {
		return nil, nil, http.StatusUnauthorized, err
	}

	meta, code, err := s.RetrieveMeta(ctx, chatID, url)
	if err != nil {
		return nil, nil, code, err
	}

	link, err = s.linksRepo.DeleteLink(ctx, chatID, url)
	if err != nil {
		return nil, nil, http.StatusNotFound, err
	}

	if err := s.filtersRepo.RemoveFiltersFromLink(ctx, link.LinkID); err != nil {
		return nil, nil, http.StatusBadRequest, err
	}

	if err := s.tagsRepo.RemoveTagsFromLink(ctx, link.LinkID); err != nil {
		return nil, nil, http.StatusBadRequest, err
	}

	return link, meta, http.StatusOK, nil
}

func (s *Service) RetrieveLinks(ctx context.Context, chatID int64) (links []scrapper_api.LinkResponse, status int, err error) {
	slog.Info("Retrieve Links: ", slog.Any("chatID", chatID))

	_, err = s.usersRepo.GetUser(ctx, chatID)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}

	repoLinks, err := s.linksRepo.GetLinksByUser(ctx, chatID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	links = make([]scrapper_api.LinkResponse, 0, len(repoLinks))

	for _, link := range repoLinks {
		meta, code, err := s.RetrieveMeta(ctx, chatID, link.URL)
		if err != nil {
			return nil, code, err
		}

		responseLink := scrapper_api.LinkResponse{
			Id:      &link.LinkID,
			Url:     &link.URL,
			Filters: &meta.Filters,
			Tags:    &meta.Tags,
		}

		slog.Info("Retrieve Links: Consutructed link: ",
			"linkID", link.LinkID,
			"url", link.URL,
			"filters", meta.Filters,
			"tags", meta.Tags,
		)

		links = append(links, responseLink)
	}

	return links, http.StatusOK, nil
}

func (s *Service) RetrieveMeta(ctx context.Context, chatID int64, url string) (meta *LinkMeta, code int, err error) {
	slog.Info("Retrieve Meta: ", slog.Any("chatID", chatID), slog.Any("url", url))

	link, err := s.linksRepo.GetLinkByURL(ctx, chatID, url)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	filters, err := s.filtersRepo.GetFiltersByLink(ctx, link.LinkID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	tags, err := s.tagsRepo.GetTagsByLink(ctx, link.LinkID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &LinkMeta{Tags: tags, Filters: filters}, http.StatusOK, nil
}
