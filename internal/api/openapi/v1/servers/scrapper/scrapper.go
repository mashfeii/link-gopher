package scrapper

import (
	"encoding/json"
	defaulterrors "errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/api"
	"github.com/es-debug/backend-academy-2024-go-template/internal/application"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

type API struct {
	deps *application.ScrapperDependencies
}

func NewAPI(deps *application.ScrapperDependencies) *API {
	return &API{
		deps: deps,
	}
}

// (DELETE /links).
func (s *API) DeleteLinks(w http.ResponseWriter, r *http.Request, params DeleteLinksParams) {
	slog.Info("Scrapper Endpoint: DeleteLinks: ", slog.Any("params", params))

	_, err := s.deps.Repo.GetUser(r.Context(), params.TgChatId)
	if err != nil {
		api.ResponseError(w, http.StatusUnauthorized, err)
		return
	}

	requestBody := struct {
		Link string `json:"link"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		api.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	link, err := s.deps.Repo.DeleteLink(r.Context(), params.TgChatId, requestBody.Link)
	if err != nil {
		api.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	linkFilters := make([]string, 0)

	for key, value := range link.Filters {
		for _, filter := range value {
			linkFilters = append(linkFilters, key+":"+filter)
		}
	}

	linkTags := make([]string, 0)
	for tag := range link.Tags {
		linkTags = append(linkTags, tag)
	}

	slog.Info("Scrapper Endpoint: DeleteLinks: ", "link", link)

	api.ResponseWithJSON(w, http.StatusOK, LinkResponse{
		Id:      &link.LinkID,
		Url:     &link.URL,
		Filters: &linkFilters,
		Tags:    &linkTags,
	})
}

// (GET /links).
func (s *API) GetLinks(w http.ResponseWriter, r *http.Request, params GetLinksParams) {
	slog.Info("Scrapper Endpoint: GetLinks: ", slog.Any("params", params))

	links, err := s.deps.Repo.GetLinks(r.Context(), params.TgChatId)
	if err != nil {
		api.ResponseError(w, http.StatusUnauthorized, err)
		return
	}

	linksResponse := make([]LinkResponse, 0, len(links))
	size := int32(len(links)) //nolint:gosec // API specification

	for _, link := range links {
		linkFilters := make([]string, 0)

		for key, value := range link.Filters {
			for _, filter := range value {
				linkFilters = append(linkFilters, key+":"+filter)
			}
		}

		linkTags := make([]string, 0)
		for tag := range link.Tags {
			linkTags = append(linkTags, tag)
		}

		slog.Info("Scrapper Endpoint: GetLinks: ", "link", link)

		linksResponse = append(linksResponse, LinkResponse{
			Id:      &link.LinkID,
			Url:     &link.URL,
			Filters: &linkFilters,
			Tags:    &linkTags,
		})
	}

	slog.Info("Scrapper Endpoint: GetLinks: ", "link quantity", len(links))

	api.ResponseWithJSON(w, http.StatusOK, ListLinksResponse{
		Links: &linksResponse,
		Size:  &size,
	})
}

// (POST /links).
func (s *API) PostLinks(w http.ResponseWriter, r *http.Request, params PostLinksParams) {
	slog.Info("Scrapper Endpoint: PostLinks: ", slog.Any("params", params))

	_, err := s.deps.Repo.GetUser(r.Context(), params.TgChatId)
	if err != nil {
		api.ResponseError(w, http.StatusUnauthorized, err)
		return
	}

	requestBody := struct {
		Link    string   `json:"link"`
		Tags    []string `json:"tags"`
		Filters []string `json:"filters"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		api.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	linkTags := make(map[string]struct{})
	for _, tag := range requestBody.Tags {
		linkTags[tag] = struct{}{}
	}

	linkFilters := make(map[string][]string)

	for _, filter := range requestBody.Filters {
		keyValue := strings.Split(filter, ":")
		if len(keyValue) != 2 {
			api.ResponseError(w, http.StatusBadRequest, errors.NewErrInvalidFilterFormat(filter))
			return
		}

		linkFilters[keyValue[0]] = append(linkFilters[keyValue[0]], keyValue[1])
	}

	link := &models.Link{
		ChatID:     params.TgChatId,
		URL:        requestBody.Link,
		Tags:       linkTags,
		Filters:    linkFilters,
		LastUpdate: time.Now(),
	}

	if err := s.deps.Repo.AddLink(r.Context(), params.TgChatId, link); err != nil {
		api.ResponseError(w, http.StatusInternalServerError, err)
		return
	}

	slog.Info("Scrapper Endpoint: PostLinks: ", "saved link", link)

	api.ResponseWithJSON(w, http.StatusOK, LinkResponse{
		Id:      &link.LinkID,
		Url:     &link.URL,
		Filters: &requestBody.Filters,
		Tags:    &requestBody.Tags,
	})
}

// (DELETE /tg-chat/{id}).
func (s *API) DeleteTgChatId(w http.ResponseWriter, r *http.Request, id int64) { //nolint:stylecheck,revive // generated method
	slog.Info("Scrapper Endpoint: DeleteTgChatId: ", slog.Any("id", id))

	if _, err := s.deps.Repo.DeleteUser(r.Context(), id); err != nil {
		api.ResponseError(w, http.StatusNotFound, err)
		return
	}

	api.ResponseWithJSON(w, http.StatusOK, id)
}

// (POST /tg-chat/{id}).
func (s *API) PostTgChatId(w http.ResponseWriter, r *http.Request, id int64) { //nolint:stylecheck,revive // generated method
	slog.Info("Scrapper Endpoint: PostTgChatId: ", slog.Any("id", id))

	if err := s.deps.Repo.CreateUser(r.Context(), id); err != nil {
		if defaulterrors.As(err, &errors.ErrUserAlreadyExists{}) {
			api.ResponseWithJSON(w, http.StatusAlreadyReported, id)
			return
		}

		api.ResponseError(w, http.StatusInternalServerError, err)

		return
	}

	api.ResponseWithJSON(w, http.StatusOK, id)
}
