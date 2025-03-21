package servers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	"github.com/es-debug/backend-academy-2024-go-template/internal/api/handlers/scrapper"
	scrapper_api "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/servers/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/application/service"
)

func NewScrapperServer(
	cfg *config.Config,
	service service.UserService,
	middlewares []scrapper_api.MiddlewareFunc,
) *http.Server {
	api := scrapper.NewAPI(service)

	mux := http.NewServeMux()
	handler := scrapper_api.HandlerWithOptions(api, scrapper_api.StdHTTPServerOptions{
		BaseRouter:  mux,
		Middlewares: middlewares,
	})

	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Serving.Host, cfg.Serving.ScrapperPort),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("Scrapper server is created",
		"Address", server.Addr)

	return server
}
