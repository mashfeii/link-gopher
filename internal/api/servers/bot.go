package servers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/servers/bot"
	"github.com/es-debug/backend-academy-2024-go-template/internal/application"
)

func NewBotServer(deps *application.BotDependencies) *http.Server {
	api := bot.NewAPI(deps)

	mux := http.NewServeMux()
	handler := bot.HandlerFromMux(api, mux)

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", deps.Config.Serving.Host, deps.Config.Serving.BotPort),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("Bot server is created",
		"Address", server.Addr)

	return server
}
