package servers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/servers/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/application"
)

func NewScrapperServer(deps *application.ScrapperDependencies) *http.Server {
	api := scrapper.NewAPI(deps)

	mux := http.NewServeMux()
	handler := scrapper.HandlerFromMux(api, mux)

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", deps.Config.Serving.Host, deps.Config.Serving.ScrapperPort),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("Scrapper server is created",
		"Address", server.Addr)

	return server
}
