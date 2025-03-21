package scrapper

import (
	"encoding/json"
	defaulterrors "errors"
	"net/http"

	"github.com/es-debug/backend-academy-2024-go-template/internal/api"
	scrapper_api "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/servers/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/application/service"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

type API struct {
	service service.UserService
}

func NewAPI(service service.UserService) *API {
	return &API{
		service: service,
	}
}

// (DELETE /links).
func (s *API) DeleteLinks(w http.ResponseWriter, r *http.Request, params scrapper_api.DeleteLinksParams) {
	requestBody := struct {
		Link string `json:"link"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		api.ResponseError(w, http.StatusBadRequest, err)

		return
	}

	link, meta, code, err := s.service.DeleteLink(r.Context(), params.TgChatId, requestBody.Link)
	if err != nil {
		api.ResponseError(w, code, err)

		return
	}

	api.ResponseWithJSON(w, http.StatusOK, scrapper_api.LinkResponse{
		Id:      &link.LinkID,
		Url:     &link.URL,
		Filters: &meta.Filters,
		Tags:    &meta.Tags,
	})
}

// (GET /links).
func (s *API) GetLinks(w http.ResponseWriter, r *http.Request, params scrapper_api.GetLinksParams) {
	linksResponse, status, err := s.service.RetrieveLinks(r.Context(), params.TgChatId)
	if err != nil {
		api.ResponseError(w, status, err)
		return
	}

	size := int32(len(linksResponse)) //nolint:gosec // specification limitation

	api.ResponseWithJSON(w, http.StatusOK, scrapper_api.ListLinksResponse{
		Links: &linksResponse,
		Size:  &size,
	})
}

// (POST /links).
func (s *API) PostLinks(w http.ResponseWriter, r *http.Request, params scrapper_api.PostLinksParams) {
	requestBody := struct {
		Link    string   `json:"link"`
		Tags    []string `json:"tags"`
		Filters []string `json:"filters"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		api.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	link, code, err := s.service.TrackLink(
		r.Context(),
		params.TgChatId,
		requestBody.Link,
		requestBody.Tags,
		requestBody.Filters,
	)
	if err != nil {
		api.ResponseError(w, code, err)
	}

	api.ResponseWithJSON(w, http.StatusOK, scrapper_api.LinkResponse{
		Id:      &link.LinkID,
		Url:     &link.URL,
		Filters: &requestBody.Filters,
		Tags:    &requestBody.Tags,
	})
}

// (DELETE /tg-chat/{id}).
func (s *API) DeleteTgChatId(w http.ResponseWriter, r *http.Request, id int64) { //nolint:stylecheck,revive // generated method
	code, err := s.service.DeleteUser(r.Context(), id)
	if err != nil {
		api.ResponseError(w, code, err)
		return
	}

	api.ResponseWithJSON(w, code, id)
}

// (POST /tg-chat/{id}).
func (s *API) PostTgChatId(w http.ResponseWriter, r *http.Request, id int64) { //nolint:stylecheck,revive // generated method
	code, err := s.service.RegisterUser(r.Context(), id)
	if err != nil {
		if defaulterrors.As(err, &errors.ErrUserAlreadyExists{}) {
			api.ResponseWithJSON(w, http.StatusAlreadyReported, id)
			return
		}

		api.ResponseError(w, code, err)

		return
	}

	api.ResponseWithJSON(w, code, id)
}
