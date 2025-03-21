package bot

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/es-debug/backend-academy-2024-go-template/internal/api"
	"github.com/es-debug/backend-academy-2024-go-template/internal/application"
)

type API struct {
	deps *application.BotDependencies
}

func NewAPI(deps *application.BotDependencies) *API {
	return &API{
		deps: deps,
	}
}

// (POST /updates).
func (s *API) PostUpdates(w http.ResponseWriter, r *http.Request) {
	slog.Info("Bot Endpoint: PostUpdates")

	requestBody := struct {
		Description string `json:"description"`
		TgChatID    int64  `json:"tgChatId"`
		URL         string `json:"url"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		api.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	botClient := s.deps.TgClient.GetBot()
	message := tgbotapi.NewMessage(requestBody.TgChatID,
		fmt.Sprintf("New update from your favorite website: %s\nUpdate: %s", requestBody.URL, requestBody.Description),
	)

	_, err := botClient.Send(message)
	if err != nil {
		api.ResponseError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
